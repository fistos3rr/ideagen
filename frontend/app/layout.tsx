import type { Metadata } from "next";
import { AuthProvider } from '@/lib/auth/AuthProvider';
import "./globals.css";

export const metadata: Metadata = {
  title: "IdeaGen App",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ru">
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
