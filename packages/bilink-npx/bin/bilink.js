#!/usr/bin/env node
import { ensureBinary, runBinary } from "../lib/download.js";
import { releaseUrl, resolveTarget } from "../lib/platform.js";
import path from "node:path";
import { homedir } from "node:os";

const { binary } = resolveTarget(process.platform, process.arch);
const cacheDir = path.join(homedir(), ".cache", "bilink");
const binPath = path.join(cacheDir, binary);

try {
  await ensureBinary(binPath, releaseUrl("latest", binary));
} catch (err) {
  console.error(`bilink download failed: ${err?.message ?? err}`);
  process.exit(1);
}

process.exit(runBinary(binPath, process.argv.slice(2)));
