'use client'

import { ArrowLeft } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import apiClient from '@/lib/apiClient'
import { computeStatus } from '@/lib/orderStatus'
import Link from 'next/link'
import PayButton from './payButton'
import CountdownTimer from './countdown'
import { useState, useEffect } from 'react'
import { FaStar } from 'react-icons/fa';
import { Input } from "@/components/ui/input"

export default function OrderPage({ params, userType = "buyer", }) {
  const [order, setOrder] = useState(null)
  const [addresses, setAddresses] = useState([])
  const [selectedAddress, setSelectedAddress] = useState(null)
  const [orderId, setOrderId] = useState(null)  
  const [reviews, setReviews] = useState([])
  const [activeReviewIndex, setActiveReviewIndex] = useState(null);
  const [coupons, setCoupons] = useState([]);
  const [selectedCoupon, setSelectedCoupon] = useState(null);
  const [discountedPrice, setDiscountedPrice] = useState(order ? order.total_price : 0);
  const [couponCode, setCouponCode] = useState('');
  const [productReviews, setProductReviews] = useState({});

  const toggleReviewInput = (index) => {
    setActiveReviewIndex(activeReviewIndex === index ? null : index);
  };

  const StarRating = ({ rating, onRatingChange }) => {
    return (
      <div className="flex">
        {[...Array(5)].map((_, index) => {
          const starValue = index + 1;
          return (
            <FaStar
              key={index}
              size={24}
              color={starValue <= rating ? "#ffc107" : "#e4e5e9"}
              onClick={() => onRatingChange(starValue)}
              style={{ cursor: "pointer" }}
            />
          );
        })}
      </div>
    );
  };

  useEffect(() => {
    // Unwrap params and fetch order details
    const fetchData = async () => {
      const { id } = await params
      setOrderId(id)

      try {
        // Fetch order details
        const orderRes = await apiClient.get(`/order/${id}`)
        const fetchedOrder = orderRes.data.order
        setOrder(fetchedOrder)

        // 获取每个商品的评论
        const reviewsPromises = fetchedOrder.items.map(item => 
          apiClient.get(`/product/${item.product_id}/reviews`)
        );
        const reviewsResponses = await Promise.all(reviewsPromises);
        const reviewsMap = {};
        reviewsResponses.forEach((response, index) => {
          const productId = fetchedOrder.items[index].product_id;
          reviewsMap[productId] = response.data.reviews;
        });
        setProductReviews(reviewsMap);

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
    if (order) {
      setReviews(order.items.map(() => ({ rating: 0, comment: '' })));
    }
  }, [order]);

  useEffect(() => {
    const fetchCoupons = async () => {
      try {
        const res = await apiClient.get('/product/promotions');
        setCoupons(res.data.promotions);
      } catch (error) {
        console.error('Error fetching coupons:', error);
      }
    };

    fetchCoupons();
  }, []);

  useEffect(() => {
    if (order && order.discount && order.discount !== 0) {
      setDiscountedPrice(order.total_price - order.discount);
    }
  }, [order]);

  const handleAddressChange = (e) => {
    const addressIndex = e.target.value
    setSelectedAddress(addresses[addressIndex])
  }

  const handleCancelOrder = async () => {
    try {
      await apiClient.put(`/order/status`, {
        order_id: parseInt(orderId),
        status: 4
      });
      alert("Order cancelled");
      window.location.reload();
    } catch (error) {
      console.error('Error cancelling order:', error);
    }
  };

  const handleConfirmReceipt = async () => {
    try {
      await apiClient.put(`/order/status`, {
        order_id: parseInt(orderId),
        status: 3
      });
      alert('Order confirmed');
      window.location.reload();
    } catch (error) {
      console.error('Error confirming receipt:', error);
    }
  };

  const handleRequestRefund = async () => {
    try {
      await apiClient.put(`/order/status`, {
        order_id: parseInt(orderId),
        status: 5
      });
      alert('Refund request submitted');
      window.location.reload();
    } catch (error) {
      console.error('Error requesting refund:', error);
    }
  };

  const handleReviewChange = (index, field, value) => {
    const newReviews = [...reviews];
    newReviews[index][field] = value;
    setReviews(newReviews);
  };

  const handleSubmitReview = async (index) => {
    const { rating, comment } = reviews[index];
    const productId = order.items[index].product_id;
    
    // 检查是否已经评论过
    const existingReviews = productReviews[productId] || [];
    const hasReviewed = existingReviews.some(review => review.user_id === order.user_id);
    
    if (hasReviewed) {
      alert('You have already commented on this product');
      return;
    }

    try {
      await apiClient.post(`/product/reviews`, {
        product_id: productId,
        user_id: order.user_id,
        rating,
        comment,
      });
      alert('Review submitted');

      // 更新本地状态
      const updatedReviews = [...reviews];
      updatedReviews[index] = { rating, comment };
      setReviews(updatedReviews);

      // 重新获取该商品的评论
      const reviewsRes = await apiClient.get(`/product/${productId}/reviews`);
      setProductReviews(prev => ({
        ...prev,
        [productId]: reviewsRes.data.reviews
      }));
    } catch (error) {
      console.error('Error submitting review:', error);
    }
  };

  const handleApplyCoupon = async () => {
    try {
      const coupon = coupons.find(c => c.coupon_code === couponCode);
      if (!coupon) {
        alert('Invalid coupon code');
        return;
      }

      let discount = 0;
      if (coupon.coupon_type === 'fixed') {
        discount = coupon.coupon_value;
      } else if (coupon.coupon_type === 'percentage') {
        discount = Math.min(order.total_price * coupon.coupon_value / 100, coupon.usage_limit);
      }

      const newPrice = order.total_price - discount;
      setDiscountedPrice(newPrice);
      setSelectedCoupon(coupon);

      // 自动更新订单金额
      await apiClient.put(`/order/discount`, {
        order_id: parseInt(orderId),
        discount: discount.toString()
      });
      window.location.reload();
      alert('Coupon applied');
    } catch (error) {
      console.error('Error applying coupon:', error);
    }
  };

  if (!order) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <Link href="/orders">
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
              {order.status != 0 ? (
              <>
                {/* // Display order's shipping address if already paid */}
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
              </>
              ):(
                // Allow user to select an address if not paid
                <div>
                  <label htmlFor="address-select" className="block mb-2">
                    Select Address:
                  </label>
                  <select
                    id="address-select"
                    className="w-full p-2 border rounded"
                    onChange={handleAddressChange}
                    value={addresses.indexOf(selectedAddress)}
                  >
                    {addresses.map((address, index) => (
                      <option key={index} value={index}>
                        {address.address_name}, {address.address_phone}, {address.address}
                      </option>
                    ))}
                  </select>
                  {selectedAddress && (
                    <div className="mt-4">
                      <div>
                        <strong>Name:</strong> {selectedAddress.address_name}
                      </div>
                      <div>
                        <strong>Phone:</strong> {selectedAddress.address_phone}
                      </div>
                      <div>
                        <strong>Address:</strong> {selectedAddress.address}
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
          {order.status === 0 && (
          <div className="mt-4">
            <h3 className="font-semibold mb-2">Enter Coupon Code</h3>
            <Input 
              value={couponCode} 
              onChange={(e) => setCouponCode(e.target.value)} 
              placeholder="Enter coupon code" 
              className="mb-2"
            />
            <Button onClick={handleApplyCoupon} disabled={!couponCode}>
              Apply Coupon
            </Button>
          </div>
          )}
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
                    <TableCell>
                      {order.status === 3 && (
                        <>
                          {(() => {
                            const productId = item.product_id;
                            const existingReviews = productReviews[productId] || [];
                            const userReview = existingReviews.find(review => review.user_id === order.user_id);
                            
                            if (userReview) {
                              return (
                                <div className="space-y-2">
                                  <div className="flex items-center">
                                    {[...Array(5)].map((_, i) => (
                                      <FaStar
                                        key={i}
                                        size={16}
                                        color={i < userReview.rating ? "#ffc107" : "#e4e5e9"}
                                      />
                                    ))}
                                  </div>
                                  <p className="text-sm text-gray-600">{userReview.comment}</p>
                                </div>
                              );
                            }
                            
                            return (
                              <Button onClick={() => toggleReviewInput(idx)} className="mt-2">
                                {activeReviewIndex === idx ? 'Hide Review' : 'Review'}
                              </Button>
                            );
                          })()}
                        </>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          {activeReviewIndex !== null && (
          <div className="mt-4">
            <h3 className="font-semibold mb-2">Rate & Review</h3>
            <div className="mb-4">
              <div className="font-medium">{order.items[activeReviewIndex].product_name}</div>
              <div className="flex items-center space-x-2 mt-2">
                <label>Rating:</label>
                <StarRating
                  rating={reviews[activeReviewIndex]?.rating || 0}
                  onRatingChange={(newRating) => handleReviewChange(activeReviewIndex, 'rating', newRating)}
                />
              </div>
              <div className="mt-2">
                <label>Comment:</label>
                <textarea
                  value={reviews[activeReviewIndex]?.comment || ''}
                  onChange={(e) => handleReviewChange(activeReviewIndex, 'comment', e.target.value)}
                  className="w-full p-2 border rounded"
                />
              </div>
              <Button onClick={() => handleSubmitReview(activeReviewIndex)} className="mt-2">
                Submit Review
              </Button>
            </div>
          </div>
        )}
          {order.status === 0 && (
            <div className="flex justify-center mt-4 space-x-8">
              <Button onClick={handleCancelOrder} variant="destructive">Cancel</Button>
              <PayButton orderID={orderId} address={selectedAddress} price={discountedPrice}>
                Pay
              </PayButton>
            </div>
          )}
          {(order.status === 1 || order.status === 2) && (
          <div className="flex justify-center mt-4 space-x-8">
            <Button onClick={handleConfirmReceipt} className="bg-green-500 hover:bg-green-600 text-white">
              Confirm Receipt
            </Button>
            <Button onClick={handleRequestRefund} className="bg-red-500 hover:bg-red-600 text-white">
              Request Refund
            </Button>
          </div>
        )}
        </CardContent>
      </Card>
    </div>
  )
}