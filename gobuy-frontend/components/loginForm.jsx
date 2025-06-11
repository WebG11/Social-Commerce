'use client'

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useState } from "react"
import apiClient from "@/lib/apiClient"

import { useRouter } from 'next/navigation'
export function LoginForm({
  className,
  userType = "buyer",
  ...props
}) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const router = useRouter()

  const themeConfig = {
    buyer: {
      primary: 'text-yellow-900',
      button: 'bg-yellow-400 hover:bg-yellow-500', //'bg-green-400 hover:bg-green-500',
      border: 'border-yellow-900'//'border-green-400'
    },
    seller: {
      primary: 'text-blue-900',
      button: 'bg-blue-400 hover:bg-blue-500',
      border: 'border-blue-900'
    }
  }

  const handleLogin = async()=>{
    try {
      const res = await apiClient.post("/login", {
        email,
        password,
        is_seller: userType === 'seller' // 根据用户类型添加字段
      })
      const {username,access_token, userId} = res.data;
      localStorage.setItem("username", username)
      localStorage.setItem("is_seller", userType === 'seller')
      localStorage.setItem("userId", userId)
      document.cookie = `Authorization=Bearer ${access_token}; Path=/; SameSite=Lax`;

      console.log( document.cookie)
      router.push("/")
    }  catch (error) {
      if (error.response && error.response.status === 401) {
        alert("Login failed: Incorrect email or password.");
      } else {
        console.error("Login error:", error);
        alert("Login failed: Incorrect email or password.");
      }
    }
  }
  
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle className={`text-2xl ${themeConfig[userType].primary}`}>
          {userType === 'seller' ? 'Sign in with business credentials' : 'Sign in'}
          </CardTitle>
          {/* <CardDescription>
            Enter your email and password below to login
          </CardDescription> */}
        </CardHeader>
        <CardContent>
          <form>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
              <div className="grid gap-2">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <a
                    href="#"
                    className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                  >
                    Forget password?
                  </a>
                </div>
                <Input id="password" value={password}
                  onChange={(e) => setPassword(e.target.value)} type="password" required />
              </div>
              <Button 
               onClick={handleLogin} 
               type="button" 
               className={`w-full text-black rounded-full ${themeConfig[userType].border} ${themeConfig[userType].button}`}>
                Continue
              </Button>

              {/* 新增虚线分割线和商家登录提示 */}
            {userType === 'buyer' && (
              <div className="mt-6 space-y-2 text-center">
                <div className="border-t border-dashed border-gray-300 my-4"></div>
                <div className="text-sm text-gray-600">
                  <p>Buying for work?</p>
                  <a 
                    href="/seller/login" 
                    className="font-medium text-blue-600 hover:text-blue-500 hover:underline"
                  >
                    Shop on Gobuy Business
                  </a>
                </div>
              </div>
            )}
            </div>
          </form>
        </CardContent>
      </Card>
      {/* 分割线 */}
      <div className="relative my-6">
        <div className="absolute inset-0 flex items-center">
          <div className={`w-full border-t ${themeConfig[userType].border}`}></div>
        </div>
        <div className="relative flex justify-center text-sm">
          <span className={`bg-white px-2 ${themeConfig[userType].primary}`}>New to Goby?</span>
        </div>
      </div>

      {/* 创建账户按钮 */}
      <Button
        variant="outline"
        className={`w-full rounded-full ${themeConfig[userType].border} ${themeConfig[userType].primary}`}
        onClick={() => router.push(userType === 'seller' ? '/seller/register' : '/register')}
      >
        Create your Goby {userType === 'seller' ? 'seller' : ''} account
      </Button>

    </div>
  )
}