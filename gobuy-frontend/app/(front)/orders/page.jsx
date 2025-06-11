'use client'

import Link from 'next/link'
import { Eye } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import apiClient from '@/lib/apiClient'
import { computeStatus } from '@/lib/orderStatus'
import { useState, useEffect } from 'react'

export default function OrdersPage() {
  const [orders, setOrders] = useState([])
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('')

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const res = await apiClient.get('/order/user'); // 调用后端 API 获取订单数据
        let fetchedOrders = res.data.orders;
        // 根据状态过滤订单
        if (statusFilter && statusFilter !== 'all') {
          fetchedOrders = fetchedOrders.filter(order => computeStatus(order.status) === statusFilter);
        }
        setOrders(fetchedOrders);
      } catch (error) {
        console.error('Error fetching orders:', error);
      }
    };

    fetchOrders();
  }, [statusFilter]); // 空依赖数组，确保只在组件挂载时运行一次

  return (
    <div className="container mx-auto">
      <h1 className="text-3xl font-bold mb-6">My Orders</h1>
      <div className="flex flex-col sm:flex-row gap-4 mb-6">
        <Input
          placeholder="Search order..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="max-w-sm"
        />
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Statuses</SelectItem>
            <SelectItem value="Pending">Pending</SelectItem>
            <SelectItem value="Paid">Paid</SelectItem>
            <SelectItem value="Shipped">Shipped</SelectItem>
            <SelectItem value="Delivered">Delivered</SelectItem>
            <SelectItem value="Cancelled">Cancelled</SelectItem>
            <SelectItem value="Refunded">Refunded</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {orders.map((order) => (
          <Card key={order.id} className="flex flex-col">
            <CardContent className="flex-grow p-6">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-sm">Order Number: {order.number}</h3>
                  <div className="text-sm">Created At: {order.created_at}</div>
                </div>
                <Badge className={"w-18 text-center"} variant={
                  order.status === 0 ? 'secondary' : order.status === 5 ? 'destructive' : 'success'
                }>
                  {computeStatus(order.status)}
                </Badge>
              </div>
              <div className="space-y-4">
                {order.items.map((item,idx) => (
                  <div key={idx} className="flex items-center space-x-4">
                    <div className="relative w-16 h-16 rounded overflow-hidden">
                      <img
                        src={item.product_image}
                        alt={item.product_name}
                        layout="fill"
                      />
                    </div>
                    <div className="flex-grow">
                      <div className="font-medium">{item.product_name}</div>
                      <div className="text-sm text-muted-foreground">
                        <div>Quantity: {item.quantity}</div>
                        <div>Price: ￥{item.price}</div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
            <CardFooter className="flex justify-between items-center border-t pt-6">
              <p className="font-semibold">Total Price: ￥{order.total_price}</p>
              <Link href={`/orders/${order.id}`}>
                <Button variant="outline" size="sm">
                  <Eye className="mr-2 h-4 w-4" />
                  View Details
                </Button>
              </Link>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  )
}

