export function resolveTarget(platform, arch) {
  const os = platform === "win32" ? "windows" : platform;
  if (os !== "darwin" && os !== "linux" && os !== "windows") {
    throw new Error(`unsupported platform: ${platform}-${arch}`);
  }
  let mappedArch = arch;
  if (arch === "x64") mappedArch = "amd64";
  if (arch !== "x64" && arch !== "arm64") {
    throw new Error(`unsupported platform: ${platform}-${arch}`);
  }
  const ext = os === "windows" ? ".exe" : "";
  return { binary: `bilink-${os}-${mappedArch}${ext}` };
}

export function releaseUrl(version, binary) {
  if (version === "latest") {
    return `https://github.com/bagakit/bilink/releases/latest/download/${binary}`;
  }
  return `https://github.com/bagakit/bilink/releases/download/${version}/${binary}`;
}
