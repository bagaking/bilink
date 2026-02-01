import assert from "node:assert/strict";
import { releaseUrl, resolveTarget } from "../lib/platform.js";

assert.equal(resolveTarget("darwin", "arm64").binary, "bilink-darwin-arm64");
assert.equal(resolveTarget("win32", "x64").binary, "bilink-windows-amd64.exe");
assert.equal(
  releaseUrl("latest", "bilink-darwin-arm64"),
  "https://github.com/bagakit/bilink/releases/latest/download/bilink-darwin-arm64"
);
