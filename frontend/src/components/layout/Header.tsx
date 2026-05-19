"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import LoginModal from "@/components/auth/LoginModal";
import { getUser, logout } from "@/lib/auth";
import { User } from "@/types/user";

export default function Header() {
  const router = useRouter();
  const [mounted, setMounted] = useState(false);

  const [user, setUser] = useState<User | null>(
    null
  );

  const [mode, setMode] = useState<"login" | "register">(
    "login"
  );

  const [open, setOpen] = useState(false);

  const [menuOpen, setMenuOpen] =
    useState(false);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setMounted(true);

    const savedUser = getUser();

    if (savedUser) {
      setUser(savedUser);
    }
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <>
      <header className="flex items-center justify-between border-b p-4">
        <Link href="/" className="text-xl font-bold hover:text-gray-600 transition-colors inline-flex items-center gap-2">
          <img src="/icon/logo.svg" alt="logo" style={{ width: 58, height: 58 }} />
          LingDU
        </Link>

        {user ? (
          <div className="flex items-center gap-2">
            <div
              className="relative"
              onMouseEnter={() => setMenuOpen(true)}
              onMouseLeave={() => setMenuOpen(false)}
            >
            <Link
              href={`/users/${user.id}`}
              className="inline-flex items-center justify-center gap-1 rounded border border-gray-300 bg-white px-4 py-2 font-medium text-gray-800 shadow-sm transition hover:bg-gray-50 hover:shadow"
            >
              <img src="/icon/me.svg" alt="user" style={{ width: 16, height: 16 }} />
              {user.username}
            </Link>

            {menuOpen && (
              <div className="absolute right-0 top-full z-10 pt-1">
                <div className="min-w-32 rounded border bg-white shadow-lg">
                  <Link
                    href={`/users/${user.id}`}
                    className="w-full px-4 py-2 text-left hover:bg-gray-100 inline-flex items-center gap-2"
                    onClick={() => setMenuOpen(false)}
                  >
                    <img src="/icon/me.svg" alt="profile" style={{ width: 16, height: 16 }} />
                    Profile
                  </Link>
                  <Link
                    href="/collections"
                    className="w-full px-4 py-2 text-left hover:bg-gray-100 inline-flex items-center gap-2"
                    onClick={() => setMenuOpen(false)}
                  >
                    <img src="/icon/collect_no.svg" alt="collections" style={{ width: 16, height: 16 }} />
                    Collections
                  </Link>
                  <Link
                    href="/history"
                    className="w-full px-4 py-2 text-left hover:bg-gray-100 inline-flex items-center gap-2"
                    onClick={() => setMenuOpen(false)}
                  >
                    <img src="/icon/history.svg" alt="history" style={{ width: 16, height: 16 }} />
                    History
                  </Link>
                  <button
                    onClick={() => {
                      logout();
                      setUser(null);
                      setMenuOpen(false);
                      router.push("/");
                      window.dispatchEvent(new Event("auth-changed"));
                    }}
                    className="w-full px-4 py-2 text-left hover:bg-gray-100 inline-flex items-center gap-2"
                  >
                    <img src="/icon/logout.svg" alt="logout" style={{ width: 16, height: 16 }} />
                    Logout
                  </button>
                </div>
              </div>
            )}
          </div>
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <button
              onClick={() => {setOpen(true); setMode("register");}}
              className="rounded border px-4 py-2 text-sm"
            >Register</button>
            <button
              onClick={() => {setOpen(true); setMode("login");}}
              className="rounded bg-black px-4 py-2 text-sm text-white"
            >Login</button>
          </div>
        )}
      </header>

      <LoginModal
        open={open}
        mode={mode}
        onClose={() => setOpen(false)}
        onSubmit={(userData) => {
          setUser(userData);
          setOpen(false);
          window.dispatchEvent(new Event("auth-changed"));
        }}
      />
    </>
  );
}