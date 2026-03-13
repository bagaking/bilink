import { spawnSync } from "node:child_process";
import { createWriteStream } from "node:fs";
import { chmod, mkdir, rename, rm, stat } from "node:fs/promises";
import http from "node:http";
import https from "node:https";
import path from "node:path";
import { pipeline } from "node:stream/promises";

const REDIRECT_CODES = new Set([301, 302, 303, 307, 308]);

export async function ensureBinary(binPath, url) {
  let needsDownload = false;
  try {
    const info = await stat(binPath);
    needsDownload = info.size === 0;
  } catch (err) {
    if (err?.code !== "ENOENT") throw err;
    needsDownload = true;
  }
  if (!needsDownload) return;
  await download(url, binPath);
  await chmod(binPath, 0o755);
}

export async function download(url, dest, options = {}) {
  const maxRedirects = options.maxRedirects ?? 5;
  await mkdir(path.dirname(dest), { recursive: true });
  const tmp = path.join(
    path.dirname(dest),
    `.${path.basename(dest)}.${process.pid}.${Date.now()}.tmp`
  );
  try {
    await downloadToTemp(url, tmp, maxRedirects);
    await rename(tmp, dest);
  } catch (err) {
    await rm(tmp, { force: true });
    throw err;
  }
}

async function downloadToTemp(url, tmp, redirectsLeft) {
  const res = await request(url);
  if (REDIRECT_CODES.has(res.statusCode)) {
    res.resume();
    if (redirectsLeft <= 0) {
      throw new Error("download failed: too many redirects");
    }
    const location = res.headers.location;
    if (!location) {
      throw new Error(`download failed: ${res.statusCode}`);
    }
    return downloadToTemp(new URL(location, url).toString(), tmp, redirectsLeft - 1);
  }
  if (res.statusCode !== 200) {
    res.resume();
    throw new Error(`download failed: ${res.statusCode}`);
  }
  await pipeline(res, createWriteStream(tmp));
}

function request(url) {
  return new Promise((resolve, reject) => {
    const client = url.startsWith("http://") ? http : https;
    const req = client.get(url, resolve);
    req.on("error", reject);
  });
}

export function runBinary(binPath, args) {
  const result = spawnSync(binPath, args, { stdio: "inherit" });
  return result.status ?? 1;
}
