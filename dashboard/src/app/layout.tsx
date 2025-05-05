import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "TastyGo Admin Dashboard",
  description: "Admin dashboard for TastyGo application",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" data-theme="tastygo">
      <body className={`${inter.className} antialiased bg-background`}>
        {children}
      </body>
    </html>
  );
}
