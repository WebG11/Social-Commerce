"use client"

import Link from 'next/link'
import { ShoppingCart, User, LogIn, LogOut,Amphora,House,TableProperties,BadgeJapaneseYen } from 'lucide-react'
import { Button } from '@/components/ui/button'
import apiClient from "@/lib/apiClient"
import { useRouter } from 'next/navigation'

// const ItemsAfterLogin = ()=>{
//   return (
    
//   )
// }



export function Header() {
  const router = useRouter()
  return (
    <header className="border-b">
      <div className="container mx-auto flex items-center justify-between py-4">
        <Link href="/" className="text-4xl font-bold">
          GoBy
        </Link>
        <nav>
          <ul className="flex items-center space-x-4">
            <li>
              <Link href="/">
                <Button variant="ghost" className="hover:text-orange-500">
                  <House className="h-5 w-5"/>
                  <span className="hover:text-orange-500">Home</span>
                  <span className="sr-only">首页</span>
                </Button>
              </Link>
            </li>
            <li>
              <Link href="/cart">
                <Button variant="ghost" className="hover:text-orange-500">
                  <ShoppingCart className="h-5 w-5" />
                  <span className="hover:text-orange-500">Cart</span>
                  <span className="sr-only">购物车</span>
                </Button>
              </Link>
            </li>
            <li>
              <Link href="/orders">
                <Button variant="ghost" className="hover:text-orange-500">
                  <BadgeJapaneseYen className="h-5 w-5" />
                  <span className="hover:text-orange-500">Order</span>
                  <span className="sr-only">订单</span>
                </Button>
              </Link>
            </li>
            <li>
              <Link href="/profile">
                <Button variant="ghost" className="hover:text-orange-500">
                  <User className="h-5 w-5" />
                  <span className="hover:text-orange-500">Profile</span>
                  <span className="sr-only">个人中心</span>
                </Button>
              </Link>
            </li>
            <li>
            <Link href="#" onClick={async (e) => {
              e.preventDefault(); // 阻止默认跳转行为
              try {
                const res = await apiClient.get('/seller'); // 向后端查询用户角色
                const { isSeller } = res.data; // 假设后端返回 { isSeller: true/false }
                if (isSeller) {
                  // 如果是 seller，跳转到后台管理页面
                  router.push('/admin/products');
                } else {
                  // 如果不是 seller，跳转到 seller 登录页面
                  router.push('/seller/login');
                }
              } catch (error) {
                console.error('Error checking user role:', error);
                alert('Failed to verify user role. Please try again.');
              }
            }}>
                <Button variant="ghost" className="hover:text-orange-500">
                  <TableProperties className="h-5 w-5"/>
                  <span>Seller</span>
                  <span className="sr-only">后台管理</span>
                </Button>
              </Link>
            </li>
            <li>
              <Link href="/login">
                <Button variant="ghost" className="hover:text-orange-500">
                  <LogIn className="h-5 w-5" />
                  <span className="hover:text-orange-500">Login</span>
                  <span className="sr-only">Login</span>
                </Button>
              </Link>
            </li>
            <li>
            <Link href="#" onClick={(e) => {
                e.preventDefault();
                // 清除本地存储
                localStorage.removeItem("username");
                localStorage.removeItem("is_seller");
                localStorage.removeItem("userId");
                // 清除认证cookie
                document.cookie = "Authorization=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;";
                // 重定向到登录页面
                router.push("/login");
              }}>
                <Button variant="ghost" className="hover:text-orange-500">
                  <LogOut className="h-5 w-5" />
                  <span className="hover:text-orange-500">Logout</span>
                  <span className="sr-only">Logout</span>
                </Button>
              </Link>
            </li>
          </ul>
        </nav>
      </div>
    </header>
  )
}

