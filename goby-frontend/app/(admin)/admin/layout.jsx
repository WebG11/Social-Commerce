'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Users, ShoppingBag,BadgeJapaneseYen,Undo2,ArrowUpNarrowWide,ChartColumn } from 'lucide-react'
import { cn } from "@/lib/utils"
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"

export default function AdminLayout({ children }) {
  const pathname = usePathname()

  return (
    <div>
      <SidebarProvider>
        <div className="flex flex-1">
          <Sidebar>
            <SidebarHeader>
              <h2 className="text-xl font-bold p-4">Seller Panel</h2>
            </SidebarHeader>
            <SidebarContent>
              <SidebarMenu>
                {/* <SidebarMenuItem>
                  <SidebarMenuButton asChild className={cn(
                    "w-full justify-start",
                    pathname === "/admin/users" && "bg-muted"
                  )}>
                    <Link href="/admin/users">
                      <Users className="mr-2 h-4 w-4" />
                      User Management
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem> */}
                <SidebarMenuItem>
                  <SidebarMenuButton asChild className={cn(
                    "w-full justify-start",
                    pathname === "/admin/products" && "bg-muted"
                  )}>
                    <Link href="/admin/products">
                      <ShoppingBag className="mr-2 h-4 w-4" />
                      Product Management
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild className={cn(
                    "w-full justify-start",
                    pathname === "/admin/orders" && "bg-muted"
                  )}>
                    <Link href="/admin/orders">
                      <BadgeJapaneseYen className="mr-2 h-4 w-4" />
                      Order Management
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild className={cn(
                    "w-full justify-start",
                    pathname === "/admin/promotions" && "bg-muted"
                  )}>
                    <Link href="/admin/promotions">
                      <ArrowUpNarrowWide className="mr-2 h-4 w-4" />
                      Promotion Management
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild className={cn(
                    "w-full justify-start",
                    pathname === "/admin/salesReport" && "bg-muted"
                  )}>
                    <Link href="/admin/salesReport">
                      <ChartColumn className="mr-2 h-4 w-4" />
                      Sales Report
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>              
                <SidebarMenuItem>
                  <SidebarMenuButton asChild className={cn(
                    "w-full justify-start",
                    pathname === "/" && "bg-muted"
                  )}>
                    <Link href="/">
                      <Undo2 className="mr-2 h-4 w-4" />
                      Return to Frontend
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarContent>
          </Sidebar>
          <div className="w-full overflow-auto">
            <header className="flex items-center h-16 px-4 border-b">
              <SidebarTrigger />
              <h1 className="ml-4 text-xl font-semibold">Dashboard</h1>
            </header>
            <main className="p-4">
              {children}
            </main>
          </div>
        </div>
      </SidebarProvider>
    </div>
    
  )
}

