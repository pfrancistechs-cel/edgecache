import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "EdgeCache | Lock-Free W-TinyLFU In-Memory Cache Engine",
  description: "High-efficiency caching engine with Window TinyLFU admission, 4-bit Count-Min Sketch frequency filtering, and lock-striped SLRU.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-[#060812] text-slate-100 antialiased selection:bg-violet-500 selection:text-white">
        {children}
      </body>
    </html>
  );
}

