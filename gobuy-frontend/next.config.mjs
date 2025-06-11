// Configuration file for Next.js
/** @type {import('next').NextConfig} */
const nextConfig = {
    transpilePackages: ['msw'],
    // api: {
    //     bodyParser: {
    //       sizeLimit: '5mb'
    //     }
    // }
    images: {
      remotePatterns: [
        {
          protocol: "http",
          hostname: "**",
        },
      ],
    },
};

export default nextConfig