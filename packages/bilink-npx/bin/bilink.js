#!/usr/bin/env node
import { resolveTarget } from "../lib/platform.js";
import { spawnSync } from "node:child_process";
import { existsSync, mkdirSync } from "node:fs";
import { homedir } from "node:os";
import path from "node:path";

const { binary } = resolveTarget(process.platform, process.arch);
const cacheDir = path.join(homedir(), ".cache", "bilink");
const binPath = path.join(cacheDir, binary);

if (!existsSync(cacheDir)) mkdirSync(cacheDir, { recursive: true });
if (!existsSync(binPath)) {
  console.error("bilink binary missing; please download releases first");
  process.exit(1);
}
const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 1);
