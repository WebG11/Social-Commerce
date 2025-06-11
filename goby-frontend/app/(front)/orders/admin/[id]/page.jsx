'use client'

import { ArrowLeft } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import apiClient from '@/lib/apiClient'
import { computeStatus } from '@/lib/orderStatus'
import Link from 'next/link'
import CountdownTimer from './countdown'
import { useState, useEffect } from 'react'
import { FaStar } from 'react-icons/fa';

export default function OrderPage({ params, userType = "buyer", }) {
  const [order, setOrder] = useState(null)
  const [discountedPrice, setDiscountedPrice] = useState(order ? order.total_price : 0);
  const [orderId, setOrderId] = useState(null)  
  const [addresses, setAddresses] = useState([])
  const [selectedAddress, setSelectedAddress] = useState(null)
  
  useEffect(() => {
    // Unwrap params and fetch order details
    const fetchData = async () => {
      const { id } = await params


      try {
        // Fetch order details
        const orderRes = await apiClient.get(`/order/${id}`)
        const fetchedOrder = orderRes.data.order
        setOrder(fetchedOrder)

        // If order is already paid, no need to fetch addresses
        if (fetchedOrder.status === 1) {
          setSelectedAddress({
            address_name: fetchedOrder.name,
            address_phone: fetchedOrder.phone,
            address: fetchedOrder.address,
          })
          return
        }

        const res = await apiClient.get(`/user`); // 调用后端 API 获取用户数据
        const { name, email, password, addresses } = res.data;
        setAddresses(addresses || []);

        // Default to the first address if available
        if (addresses.length > 0) {
          setSelectedAddress(addresses[0])
        }
        
      } catch (error) {
        console.error('Error fetching data:', error)
      }
    }

    fetchData()
  }, [params])



  useEffect(() => {
    if (order && order.discount && order.discount !== 0) {
      setDiscountedPrice(order.total_price - order.discount);
    }
  }, [order]);

  if (!order) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <Link href="/admin/orders">
        <Button variant="ghost" className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          View All Orders
        </Button>
      </Link>
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Order Details</CardTitle>
          {order.status === 0 && <CountdownTimer targetTime={order.created_at} />}
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <h3 className="font-semibold mb-2">Order Information</h3>
              <div>Seller: {order.seller_id}</div>
              <div>Order Number: {order.number}</div>
              <div>Creation Date: {order.created_at}</div>
              <div>
                <h3 className="font-semibold mb-2">Total Price</h3>
                {order.discount && order.discount !== 0 ? (
                  <div>
                    <p className="text-2xl font-bold">￥{discountedPrice.toFixed(2)}</p>
                    <p className="text-sm text-gray-500">Discount Applied: ￥{order.discount.toFixed(2)}</p>
                  </div>
                ) : (
                  <p className="text-2xl font-bold">￥{order.total_price.toFixed(2)}</p>
                )}
              </div>
              <div>
                Status:{' '}
                <Badge
                  variant={
                    order.status === 0
                      ? 'secondary'// 未支付
                      : order.status === 1
                      ? 'success' // 已支付
                      : order.status === 2
                      ? 'success' // 已发货
                      : order.status === 3
                      ? 'success' // 已送达
                      : 'destructive' // 退款&取消
                  }
                >
                  {computeStatus(order.status)}
                </Badge>
              </div>
              {order.status === 2 && <div>Payment Method: Alipay</div>}
            </div>
            <div>
              <h3 className="font-semibold mb-2">Shipping Information</h3>
                <div className="mt-4">
                  <div>
                    <strong>Name:</strong> {order.name}
                  </div>
                  <div>
                    <strong>Phone:</strong> {order.phone}
                  </div>
                  <div>
                    <strong>Address:</strong> {order.address}
                  </div>
                </div>
            </div>
          </div>
          <div className="mt-4">
            <h3 className="font-semibold mb-2">Order Items</h3>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product Name</TableHead>
                  <TableHead>Image</TableHead>
                  <TableHead>Unit Price</TableHead>
                  <TableHead>Quantity</TableHead>
                  <TableHead>Total Price</TableHead>
                  {order.status === 3 && <TableHead>Review</TableHead>}
                </TableRow>
              </TableHeader>
              <TableBody>
                {order.items.map((item, idx) => (
                  <TableRow key={idx}>
                    <TableCell>{item.product_name}</TableCell>
                    <TableCell>
                      <img className="h-10 w-10" src={item.product_image} alt={item.product_name} />
                    </TableCell>
                    <TableCell>￥{item.price.toFixed(2)}</TableCell>
                    <TableCell>{item.quantity}</TableCell>
                    <TableCell>￥{(item.quantity * item.price).toFixed(2)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}