'use client';

import { useState } from "react";
import { updateUsername, changePassword } from "@/lib/api";
import { useToast } from "@/components/ui/ToastContext";
import { User } from "@/types/user";

interface Props {
  open: boolean;
  onClose: () => void;
  user: User;
  onUserUpdated: (user: User) => void;
}

export default function ProfileModal({ open, onClose, user, onUserUpdated }: Props) {
  const [tab, setTab] = useState<"username" | "password">("username");
  const [newUsername, setNewUsername] = useState(user.username);
  const [savingUsername, setSavingUsername] = useState(false);

  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [savingPassword, setSavingPassword] = useState(false);

  const { addToast } = useToast();

  if (!open) return null;

  const handleSaveUsername = async () => {
    if (!newUsername.trim()) {
      addToast("Username cannot be empty.", { type: "warning", title: "Missing field" });
      return;
    }
    setSavingUsername(true);
    try {
      const res = await updateUsername(newUsername.trim());
      if (res.code === 200) {
        onUserUpdated({ ...user, username: newUsername.trim() });
        addToast("Username updated!", { type: "success", title: "Success" });
        onClose();
      } else {
        addToast(res.message || "Failed to update.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setSavingUsername(false);
    }
  };

  const handleSavePassword = async () => {
    if (!oldPassword) {
      addToast("Please enter your current password.", { type: "warning", title: "Missing field" });
      return;
    }
    if (newPassword.length < 6) {
      addToast("New password must be at least 6 characters.", { type: "warning", title: "Too short" });
      return;
    }
    if (newPassword !== confirmPassword) {
      addToast("Passwords do not match.", { type: "warning", title: "Mismatch" });
      return;
    }
    setSavingPassword(true);
    try {
      const res = await changePassword(oldPassword, newPassword);
      if (res.code === 200) {
        addToast("Password changed!", { type: "success", title: "Success" });
        onClose();
      } else {
        addToast(res.message || "Failed to change password.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setSavingPassword(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md p-6" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-lg font-bold mb-4">Edit Profile</h2>

        {/* Tab bar */}
        <div className="flex gap-1 mb-4 border-b">
          <button
            onClick={() => setTab("username")}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition ${
              tab === "username" ? "border-black text-black" : "border-transparent text-gray-500 hover:text-gray-700"
            }`}
          >
            Username
          </button>
          <button
            onClick={() => setTab("password")}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition ${
              tab === "password" ? "border-black text-black" : "border-transparent text-gray-500 hover:text-gray-700"
            }`}
          >
            Password
          </button>
        </div>

        {/* Username form */}
        {tab === "username" && (
          <div className="space-y-3">
            <label className="block text-sm text-gray-600">New Username</label>
            <input
              className="w-full rounded border border-gray-300 p-3 text-sm"
              value={newUsername}
              onChange={(e) => setNewUsername(e.target.value)}
              disabled={savingUsername}
              placeholder={user.username}
            />
            <div className="flex justify-end gap-2">
              <button onClick={onClose} className="rounded border px-4 py-2 text-sm text-gray-600 hover:bg-gray-50">
                Cancel
              </button>
              <button
                onClick={handleSaveUsername}
                disabled={savingUsername}
                className="rounded bg-black px-4 py-2 text-sm text-white hover:bg-gray-800 disabled:opacity-50"
              >
                {savingUsername ? "Saving..." : "Save"}
              </button>
            </div>
          </div>
        )}

        {/* Password form */}
        {tab === "password" && (
          <div className="space-y-3">
            <label className="block text-sm text-gray-600">Current Password</label>
            <input
              type="password"
              className="w-full rounded border border-gray-300 p-3 text-sm"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              disabled={savingPassword}
              placeholder="Enter current password"
            />
            <label className="block text-sm text-gray-600">New Password</label>
            <input
              type="password"
              className="w-full rounded border border-gray-300 p-3 text-sm"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              disabled={savingPassword}
              placeholder="At least 6 characters"
            />
            <label className="block text-sm text-gray-600">Confirm New Password</label>
            <input
              type="password"
              className="w-full rounded border border-gray-300 p-3 text-sm"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              disabled={savingPassword}
              placeholder="Re-enter new password"
            />
            <div className="flex justify-end gap-2">
              <button onClick={onClose} className="rounded border px-4 py-2 text-sm text-gray-600 hover:bg-gray-50">
                Cancel
              </button>
              <button
                onClick={handleSavePassword}
                disabled={savingPassword}
                className="rounded bg-black px-4 py-2 text-sm text-white hover:bg-gray-800 disabled:opacity-50"
              >
                {savingPassword ? "Saving..." : "Save"}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
