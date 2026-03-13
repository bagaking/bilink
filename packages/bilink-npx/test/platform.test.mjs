import assert from "node:assert/strict";
import test from "node:test";
import { releaseUrl, resolveTarget } from "../lib/platform.js";

test("resolveTarget maps supported node platforms to release assets", () => {
  assert.equal(resolveTarget("darwin", "arm64").binary, "bilink-darwin-arm64");
  assert.equal(resolveTarget("win32", "x64").binary, "bilink-windows-amd64.exe");
});

test("resolveTarget rejects unsupported platform targets", () => {
  assert.throws(() => resolveTarget("freebsd", "x64"), new RegExp("unsupported platform: freebsd-x64"));
  assert.throws(() => resolveTarget("linux", "ia32"), new RegExp("unsupported platform: linux-ia32"));
});

test("releaseUrl resolves latest release assets", () => {
  assert.equal(
    releaseUrl("latest", "bilink-darwin-arm64"),
    "https://github.com/bagakit/bilink/releases/latest/download/bilink-darwin-arm64"
  );
});
