import fs from "node:fs/promises";
import path from "node:path";

const DOCS_ROOT = path.resolve("docs");
const OUTPUT_PATH = path.join(DOCS_ROOT, "must-sop.md");

async function listMarkdownFiles(dirPath) {
  const entries = await fs.readdir(dirPath, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    const fullPath = path.join(dirPath, entry.name);
    if (entry.isDirectory()) {
      files.push(...(await listMarkdownFiles(fullPath)));
      continue;
    }
    if (entry.isFile() && entry.name.endsWith(".md")) {
      files.push(fullPath);
    }
  }
  return files;
}

function parseFrontmatter(content) {
  if (!content.startsWith("---\n")) {
    return null;
  }
  const endIdx = content.indexOf("\n---", 4);
  if (endIdx === -1) {
    return null;
  }
  const fmRaw = content.slice(4, endIdx);
  const lines = fmRaw.split("\n");
  const data = {};
  let currentKey = null;

  for (const line of lines) {
    if (!line.trim()) {
      continue;
    }
    const keyMatch = line.match(/^([A-Za-z0-9_-]+):\s*(.*)$/);
    if (keyMatch) {
      const key = keyMatch[1];
      const value = keyMatch[2];
      currentKey = key;
      if (!value) {
        data[key] = [];
      } else {
        data[key] = value.replace(/^['"]|['"]$/g, "").trim();
      }
      continue;
    }
    if (currentKey && Array.isArray(data[currentKey])) {
      const itemMatch = line.match(/^\s*-\s*(.*)$/);
      if (itemMatch) {
        data[currentKey].push(itemMatch[1].trim());
      }
    }
  }

  return data;
}

function buildSop(entries) {
  const header = [
    "# Project SOP",
    "",
    "This SOP is generated from docs frontmatter. Do not edit manually.",
    "",
    "## Update Requirements",
    "- When a document with SOP frontmatter changes, regenerate this file with `node scripts/generate-sop.mjs` and commit the result.",
    "- Add new SOP items by updating the `sop` list in the source document frontmatter.",
    "- Keep SOP items small and actionable; use the source document for details.",
    "",
    "## SOP Items",
    ""
  ].join("\n");

  if (entries.length === 0) {
    return `${header}No SOP items found.\n`;
  }

  const sections = entries.map((entry) => {
    const lines = [
      `### ${entry.title}`,
      `Source: \`${entry.relativePath}\``
    ];
    for (const item of entry.sop) {
      lines.push(`- ${item}`);
    }
    return `${lines.join("\n")}\n`;
  });

  return `${header}${sections.join("\n")}`;
}

async function main() {
  const allFiles = await listMarkdownFiles(DOCS_ROOT);
  const entries = [];

  for (const filePath of allFiles) {
    if (path.basename(filePath) === "must-sop.md") {
      continue;
    }
    const content = await fs.readFile(filePath, "utf8");
    const frontmatter = parseFrontmatter(content);
    if (!frontmatter || !frontmatter.title || !Array.isArray(frontmatter.sop)) {
      continue;
    }
    entries.push({
      title: frontmatter.title,
      sop: frontmatter.sop,
      relativePath: path.relative(process.cwd(), filePath)
    });
  }

  entries.sort((a, b) => a.relativePath.localeCompare(b.relativePath));
  const output = buildSop(entries);
  await fs.writeFile(OUTPUT_PATH, output, "utf8");
}

await main();
