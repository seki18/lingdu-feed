"use client";

import { useEffect, useRef, useState } from "react";
import { apiFetch, ApiResponse } from "@/lib/api";
import { saveAuth } from "@/lib/auth";
import { useToast } from "@/components/ui/ToastContext";
import { User } from "@/types/user";

interface Props {
  open: boolean;
  mode: string;
  onClose: () => void;
  onSubmit: (user: User) => void;
}

export default function LoginModal({
  open,
  mode,
  onClose,
  onSubmit,
}: Props) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(false);
  const { addToast } = useToast();

  // When switching to register mode, forcibly clear any autofilled values
  const prevMode = useRef(mode);
  useEffect(() => {
    if (mode === "register" && prevMode.current !== mode) {
      // Use a microtask to let browser autofill settle, then clear
      const t = setTimeout(() => {
        setEmail("");
        setPassword("");
        setConfirmPassword("");
      }, 50);
      prevMode.current = mode;
      return () => clearTimeout(t);
    }
    prevMode.current = mode;
  }, [mode]);

  if (!open) {
    return null;
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  const handleSubmit = async () => {
    if (!email.trim() || !password.trim()) {
      addToast("Email and password are required", {
        type: "warning",
        title: "Missing fields",
      });
      return;
    }

    if (!emailRegex.test(email.trim())) {
      addToast("Please enter a valid email address", {
        type: "warning",
        title: "Invalid email",
      });
      return;
    }

    if (mode === "register") {
      if (!username.trim()) {
        addToast("Username is required for registration", {
          type: "warning",
          title: "Missing fields",
        });
        return;
      }
      if (password !== confirmPassword) {
        addToast("Passwords do not match", {
          type: "warning",
          title: "Password mismatch",
        });
        return;
      }
    }

    setLoading(true);
    try {
      const response: ApiResponse = await apiFetch("/auth/" + mode, {
        method: "POST",
        body: JSON.stringify({
          email,
          password,
          ...(mode === "register" && {
            username,
          }),
        }),
      });

      if (response.code !== 200) {
        const errorMsg = response.message;
        addToast(errorMsg, {
          type: "error",
          title: mode === "login" ? "Login failed" : "Registration failed",
        });
        return;
      }

      const data = response.data;
      saveAuth(data.token, data.user);
      addToast(
        mode === "login"
          ? "Logged in successfully"
          : "Account created successfully",
        {
          type: "success",
          title: mode === "login" ? "Login" : "Registration",
        }
      );
      setEmail("");
      setPassword("");
      setConfirmPassword("");
      setUsername("");
      onSubmit(data.user);
    } catch (error) {
      console.error(error);
      addToast("An unexpected error occurred", {
        type: "error",
        title: "Error",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/50">
      <div className="w-100 rounded bg-white p-6 shadow-lg">
        <h2 className="mb-4 text-xl font-bold">
          {mode === "login" ? "Login" : "Register"}
        </h2>

        <div className="flex flex-col gap-4">
          <input
            className="border p-2"
            placeholder="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={loading}
            autoComplete={mode === "login" ? "email" : "off"}
            name={mode === "login" ? "email" : "reg-email"}
          />

          <input
            className="border p-2"
            type="password"
            placeholder="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={loading}
            autoComplete={mode === "login" ? "current-password" : "off"}
            name={mode === "login" ? "password" : "reg-password"}
          />

          {mode === "register" ? (
            <>
              <input
                className="border p-2"
                type="password"
                placeholder="confirm password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                disabled={loading}
                autoComplete="off"
                name="reg-confirm-password"
              />
              <input
                className="border p-2"
                placeholder="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={loading}
                autoComplete="off"
                name="reg-username"
              />
            </>
          ) : (
            <></>
          )}

          <button
            className="rounded bg-black py-2 text-white disabled:cursor-not-allowed disabled:opacity-50"
            onClick={handleSubmit}
            disabled={loading}
          >
            {loading
              ? mode === "login"
                ? "Logging in..."
                : "Registering..."
              : mode === "login"
                ? "Login"
                : "Register"}
          </button>

          <button
            className="border py-2 disabled:cursor-not-allowed disabled:opacity-50"
            onClick={onClose}
            disabled={loading}
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}