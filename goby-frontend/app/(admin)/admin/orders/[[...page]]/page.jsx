"use client"

import Link from 'next/link'
import { Eye } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import Pagination from '@/components/pagination'
import apiClient from '@/lib/apiClient'
import { computeStatus } from '@/lib/orderStatus'
import { use, useState, useEffect } from 'react'
import { ShipmentDialog } from '@/components/ui/shipmentdialog'


export default function OrderManagement({searchParams,params}){
  // const [searchTerm, setSearchTerm] = useState('')
  // const [statusFilter, setStatusFilter] = useState('')
  // const [currentPage, setCurrentPage] = useState(1)
  const [selectedOrder, setSelectedOrder] = useState(null)
  const [dialogOpen, setDialogOpen] = useState(false)

    // 发货处理函数
    const handleShipOrder = async (orderId, trackingNumber) => {
      try {
        await apiClient.put(`/admin/order/${orderId}/${trackingNumber}`)
        window.location.reload() // 简单刷新页面
      } catch (error) {
        console.error('发货失败:', error)
      }
    }

  const ordersPerPage = 10;
  const {page="1"} = use(params);
  const currentPage = Number(page);

  const [orders, setOrders] = useState([])
  const [totalPage, setTotalPage] = useState(0)

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const res = await apiClient.get("/admin/order/list", { page: currentPage, size: ordersPerPage })
        const { orders, total_count } = res.data
        setOrders(orders)
        setTotalPage(Math.ceil(total_count / ordersPerPage))
      } catch (error) {
        console.error("获取订单失败:", error)
      }
    }

    fetchOrders()
  }, [currentPage])

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Order Management</h2>
      <div className="flex gap-4 mb-4">
        {/* <Input
          placeholder="Search orders..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="max-w-sm"
        /> */}
        {/* <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Statuses</SelectItem>
            <SelectItem value="Pending">Pending</SelectItem>
            <SelectItem value="Processing">Processing</SelectItem>
            <SelectItem value="Shipped">Shipped</SelectItem>
            <SelectItem value="Delivered">Delivered</SelectItem>
          </SelectContent>
        </Select> */}
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Order Number</TableHead>
            <TableHead>Customer Name</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Total Price</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Tracking Number</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {orders.map((order) => (
            <TableRow key={order.id}>
              <TableCell>{order.number}</TableCell>
              <TableCell>{order.name}</TableCell>
              <TableCell>{order.created_at}</TableCell>
              <TableCell>￥{order.total_price-order.discount}</TableCell>
              <TableCell>
                <Badge variant={
                   order.status === 0 ? 'secondary' : (order.status === 1 || order.status === 2 || order.status === 3) ? 'success' : 'destructive'
                }>
                  {computeStatus(order.status)}
                </Badge>
              </TableCell>
              <TableCell>
                {order.tracking_number && (order.status === 2 || order.status === 3) && (
                  <div>
                    {order.tracking_number}
                  </div>
                )}
                {order.status === 1 && (
                  <Button variant="outline" size="sm" className="text-blue-900 border-blue-500 hover:bg-blue-500 hover:text-white" onClick={() => {
                    setSelectedOrder(order)
                    setDialogOpen(true)
                  }}>
                    Ship
                  </Button>
                )}
              </TableCell>
              <TableCell>
                <Link href={`/orders/admin/${order.id}`}>
                  <Button variant="ghost" size="sm">
                    <Eye className="mr-2 h-4 w-4" />
                    View Details
                  </Button>
                </Link>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Pagination hrefPrefix={'/admin/orders'} currentPage={currentPage} pageSize={ordersPerPage} totalPage={totalPage}/>
      <ShipmentDialog
        order={selectedOrder}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        onConfirm={handleShipOrder}
      />
    </div>
  )
}

