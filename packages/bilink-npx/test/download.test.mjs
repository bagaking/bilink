import assert from "node:assert/strict";
import test from "node:test";
import { createServer } from "node:http";
import { mkdtemp, readFile, rm, stat, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { download, ensureBinary } from "../lib/download.js";

async function withServer(handler, fn) {
  const server = createServer(handler);
  await new Promise((resolve) => server.listen(0, "127.0.0.1", resolve));
  try {
    const { port } = server.address();
    return await fn(`http://127.0.0.1:${port}`);
  } finally {
    await new Promise((resolve, reject) => {
      server.close((err) => (err ? reject(err) : resolve()));
    });
  }
}

async function tempPath(name) {
  const dir = await mkdtemp(path.join(tmpdir(), "bilink-npx-"));
  return { dir, file: path.join(dir, name) };
}

test("download follows redirects and writes final body", async () => {
  const { dir, file } = await tempPath("bilink");
  await withServer((req, res) => {
    if (req.url === new URL("redirect", "http://test.invalid").pathname) {
      res.writeHead(302, { location: new URL("binary", "http://test.invalid").pathname });
      res.end();
      return;
    }
    res.writeHead(200);
    res.end("binary-body");
  }, async (baseUrl) => {
    await download(new URL("redirect", `${baseUrl}/`).toString(), file);
  });
  assert.equal(await readFile(file, "utf8"), "binary-body");
  await rm(dir, { recursive: true, force: true });
});

test("download failure removes temporary file and leaves no empty cache", async () => {
  const { dir, file } = await tempPath("bilink");
  await assert.rejects(
    withServer((_req, res) => {
      res.writeHead(500);
      res.end("nope");
    }, async (baseUrl) => {
      await download(new URL("binary", `${baseUrl}/`).toString(), file);
    }),
    new RegExp("download failed: 500")
  );
  await assert.rejects(stat(file), { code: "ENOENT" });
  await rm(dir, { recursive: true, force: true });
});

test("ensureBinary redownloads zero byte cache file", async () => {
  const { dir, file } = await tempPath("bilink");
  await writeFile(file, "");
  await withServer((_req, res) => {
    res.writeHead(200);
    res.end("fresh");
  }, async (baseUrl) => {
    await ensureBinary(file, new URL("binary", `${baseUrl}/`).toString());
  });
  assert.equal(await readFile(file, "utf8"), "fresh");
  await rm(dir, { recursive: true, force: true });
});
