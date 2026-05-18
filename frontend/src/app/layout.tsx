import "./globals.css";
import Header from "@/components/layout/Header";
import { ToastProvider } from "@/components/ui/ToastContext";
import { ToastContainer } from "@/components/ui/ToastContainer";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <ToastProvider>
          <Header />
          {children}
          <ToastContainer />
        </ToastProvider>
      </body>
    </html>
  );
}