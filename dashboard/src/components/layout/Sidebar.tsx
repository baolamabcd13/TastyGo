'use client';

import { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { usePathname, useRouter } from 'next/navigation';
import { 
  FiHome, 
  FiUsers, 
  FiSettings, 
  FiLogOut,
  FiMenu,
  FiX
} from 'react-icons/fi';

export default function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const handleLogout = () => {
    // Xóa token khỏi localStorage
    localStorage.removeItem('token');
    // Chuyển hướng về trang đăng nhập
    router.push('/login');
  };

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  const menuItems = [
    {
      title: 'Trang chủ',
      icon: <FiHome className="w-5 h-5" />,
      path: '/dashboard',
    },
    {
      title: 'Quản lý Admin',
      icon: <FiUsers className="w-5 h-5" />,
      path: '/dashboard/admins',
    },
    {
      title: 'Nhật ký hoạt động',
      icon: <FiSettings className="w-5 h-5" />,
      path: '/dashboard/logs',
    },
  ];

  return (
    <>
      {/* Mobile menu button */}
      <div className="lg:hidden fixed top-4 left-4 z-50">
        <button 
          onClick={toggleMobileMenu}
          className="btn btn-circle btn-primary text-white"
        >
          {isMobileMenuOpen ? <FiX size={24} /> : <FiMenu size={24} />}
        </button>
      </div>

      {/* Sidebar for desktop */}
      <div className="hidden lg:flex flex-col h-screen w-64 bg-secondary text-white fixed">
        <div className="p-4 flex justify-center">
          <div className="relative w-32 h-16">
            <Image
              src="/tastygo.png"
              alt="TastyGo Logo"
              fill
              style={{ objectFit: 'contain' }}
            />
          </div>
        </div>
        
        <div className="flex-1 px-4 py-6">
          <ul className="menu menu-lg p-0 [&_li>*]:rounded-lg">
            {menuItems.map((item) => (
              <li key={item.path}>
                <Link 
                  href={item.path}
                  className={pathname === item.path ? 'active bg-primary text-white' : ''}
                >
                  {item.icon}
                  {item.title}
                </Link>
              </li>
            ))}
          </ul>
        </div>
        
        <div className="p-4 border-t border-gray-700">
          <button 
            onClick={handleLogout}
            className="btn btn-outline btn-block text-white"
          >
            <FiLogOut className="w-5 h-5" />
            Đăng xuất
          </button>
        </div>
      </div>

      {/* Mobile sidebar */}
      {isMobileMenuOpen && (
        <div className="lg:hidden fixed inset-0 z-40 bg-black bg-opacity-50">
          <div className="bg-secondary text-white h-screen w-64 p-4 flex flex-col">
            <div className="flex justify-between items-center mb-6">
              <div className="relative w-32 h-16">
                <Image
                  src="/tastygo.png"
                  alt="TastyGo Logo"
                  fill
                  style={{ objectFit: 'contain' }}
                />
              </div>
              <button 
                onClick={toggleMobileMenu}
                className="btn btn-circle btn-sm"
              >
                <FiX size={20} />
              </button>
            </div>
            
            <ul className="menu menu-lg p-0 [&_li>*]:rounded-lg">
              {menuItems.map((item) => (
                <li key={item.path}>
                  <Link 
                    href={item.path}
                    className={pathname === item.path ? 'active bg-primary text-white' : ''}
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    {item.icon}
                    {item.title}
                  </Link>
                </li>
              ))}
            </ul>
            
            <div className="mt-auto p-4 border-t border-gray-700">
              <button 
                onClick={handleLogout}
                className="btn btn-outline btn-block text-white"
              >
                <FiLogOut className="w-5 h-5" />
                Đăng xuất
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
