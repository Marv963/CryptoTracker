import Providers from "@/components/Providers";
import type { Metadata } from "next";
import localFont from "next/font/local";
import { Toaster } from "sonner";

import "./globals.css";
import { default as Navbar } from "@/components/navbar/Navbar";

const fonts = localFont({
  src: [
    {
      path: "../../public/fonts/CalSans-SemiBold.otf",
      weight: "600",
      style: "normal",
    },
  ],
});

export const metadata: Metadata = {
  title: "CryptoTracker",
  description: "CryptoTracker",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={`${fonts.className}`}>
        <Providers>
          <Navbar />
          <main className="container pt-12">{children}</main>
          <Toaster />
        </Providers>
      </body>
    </html>
  );
}
