import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Use standalone for Docker, default for Netlify
  output: process.env.NETLIFY ? undefined : "standalone",
};

export default nextConfig;
