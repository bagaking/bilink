#!/usr/bin/env node
import { releaseUrl, resolveTarget } from "../lib/platform.js";
import { spawnSync } from "node:child_process";
import { chmodSync, createWriteStream, existsSync, mkdirSync } from "node:fs";
import { homedir } from "node:os";
import path from "node:path";
import https from "node:https";

const { binary } = resolveTarget(process.platform, process.arch);
const cacheDir = path.join(homedir(), ".cache", "bilink");
const binPath = path.join(cacheDir, binary);

if (!existsSync(cacheDir)) mkdirSync(cacheDir, { recursive: true });

if (!existsSync(binPath)) {
  try {
    await download(releaseUrl("latest", binary), binPath);
    chmodSync(binPath, 0o755);
  } catch (err) {
    console.error(`bilink download failed: ${err?.message ?? err}`);
    process.exit(1);
  }
}

const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 1);

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest);
    https
      .get(url, (res) => {
        if (res.statusCode !== 200) {
          reject(new Error(`download failed: ${res.statusCode}`));
          return;
        }
        res.pipe(file);
        file.on("finish", () => file.close(resolve));
      })
      .on("error", reject);
  });
}
