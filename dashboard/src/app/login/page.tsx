"use client";

import { useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import axios from "axios";
import { FiLock, FiMail, FiAlertCircle } from "react-icons/fi";

type LoginFormData = {
  email: string;
  password: string;
};

export default function LoginPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>();

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    setError("");

    try {
      const response = await axios.post(
        "http://localhost:8080/api/auth/login",
        data
      );

      if (response.data && response.data.token) {
        // Lưu token vào localStorage
        localStorage.setItem("token", response.data.token);

        // Chuyển hướng đến trang dashboard
        router.push("/dashboard");
      } else {
        setError("Đăng nhập không thành công");
      }
    } catch (err: any) {
      console.error("Login error:", err);
      setError(err.response?.data?.error || "Đã xảy ra lỗi khi đăng nhập");
    } finally {
      setIsLoading(false);
    }
  };

  const [showForgotPassword, setShowForgotPassword] = useState(false);
  const [forgotEmail, setForgotEmail] = useState("");
  const [forgotEmailSent, setForgotEmailSent] = useState(false);

  const handleForgotPassword = (e: React.FormEvent) => {
    e.preventDefault();
    // Giả lập gửi email khôi phục mật khẩu
    setForgotEmailSent(true);
    setTimeout(() => {
      setShowForgotPassword(false);
      setForgotEmailSent(false);
    }, 3000);
  };

  return (
    <div className="min-h-screen flex items-center justify-center login-background p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl backdrop-blur-sm bg-opacity-95 relative z-10">
        <div className="card-body p-8">
          <div className="flex justify-center mb-6">
            <div className="w-32 h-32 relative">
              <Image
                src="/tastygo.png"
                alt="TastyGo Logo"
                fill
                style={{ objectFit: "contain" }}
                priority
              />
            </div>
          </div>

          <p className="text-center text-gray-600 mb-8">
            Đăng nhập để quản lý hệ thống
          </p>

          {error && (
            <div className="alert alert-error mb-6 shadow-md">
              <FiAlertCircle className="w-6 h-6" />
              <span>{error}</span>
            </div>
          )}

          {!showForgotPassword ? (
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium flex items-center gap-2">
                    <FiMail className="text-primary" /> Email
                  </span>
                </label>
                <div className="relative">
                  <input
                    type="email"
                    placeholder="Nhập email của bạn"
                    className={`input input-bordered w-full pl-4 ${
                      errors.email ? "input-error" : "focus:border-primary"
                    }`}
                    {...register("email", {
                      required: "Email là bắt buộc",
                      pattern: {
                        value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                        message: "Email không hợp lệ",
                      },
                    })}
                  />
                </div>
                {errors.email && (
                  <label className="label">
                    <span className="label-text-alt text-error flex items-center gap-1">
                      <FiAlertCircle className="w-3 h-3" />{" "}
                      {errors.email.message}
                    </span>
                  </label>
                )}
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium flex items-center gap-2">
                    <FiLock className="text-primary" /> Mật khẩu
                  </span>
                </label>
                <div className="relative">
                  <input
                    type="password"
                    placeholder="Nhập mật khẩu của bạn"
                    className={`input input-bordered w-full pl-4 ${
                      errors.password ? "input-error" : "focus:border-primary"
                    }`}
                    {...register("password", {
                      required: "Mật khẩu là bắt buộc",
                      minLength: {
                        value: 6,
                        message: "Mật khẩu phải có ít nhất 6 ký tự",
                      },
                    })}
                  />
                </div>
                {errors.password && (
                  <label className="label">
                    <span className="label-text-alt text-error flex items-center gap-1">
                      <FiAlertCircle className="w-3 h-3" />{" "}
                      {errors.password.message}
                    </span>
                  </label>
                )}
              </div>

              <div className="flex justify-end">
                <button
                  type="button"
                  className="text-sm text-primary hover:underline"
                  onClick={() => setShowForgotPassword(true)}
                >
                  Quên mật khẩu?
                </button>
              </div>

              <div className="form-control mt-2">
                <button
                  type="submit"
                  className={`btn btn-primary text-white ${
                    isLoading ? "loading" : ""
                  }`}
                  disabled={isLoading}
                >
                  {isLoading ? "Đang đăng nhập..." : "Đăng nhập"}
                </button>
              </div>
            </form>
          ) : (
            <div className="space-y-6">
              <h3 className="text-xl font-semibold text-center mb-4">
                Khôi phục mật khẩu
              </h3>

              {forgotEmailSent ? (
                <div className="alert alert-success">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="stroke-current shrink-0 h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>Đã gửi email khôi phục mật khẩu!</span>
                </div>
              ) : (
                <form onSubmit={handleForgotPassword}>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Email</span>
                    </label>
                    <input
                      type="email"
                      placeholder="Nhập email của bạn"
                      className="input input-bordered w-full"
                      value={forgotEmail}
                      onChange={(e) => setForgotEmail(e.target.value)}
                      required
                    />
                  </div>

                  <div className="form-control mt-6">
                    <button
                      type="submit"
                      className="btn btn-primary text-white"
                    >
                      Gửi email khôi phục
                    </button>
                  </div>
                </form>
              )}

              <div className="text-center mt-4">
                <button
                  className="text-sm text-primary hover:underline"
                  onClick={() => setShowForgotPassword(false)}
                >
                  Quay lại đăng nhập
                </button>
              </div>
            </div>
          )}

          <div className="divider mt-8"></div>
        </div>
      </div>
    </div>
  );
}
