"use client"
import { useEffect, useState, use } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Loader2, CheckCircle2, Shield } from 'lucide-react'
import { QRCode } from 'react-qr-code'
import apiClient from '@/lib/apiClient'

const paymentMethods = [
  {
    id: "alipay",
    name: "Alipay",
    icon: "/images/alipay.png",
    color: "#1677FF",
    type: "qrcode"
  },
  {
    id: "wechat",
    name: "Wechat Pay",
    icon: "/images/wechatpay.png",
    color: "#07C160",
    type: "qrcode"
  },
  {
    id: "bank",
    name: "Credit Card",
    icon: "/images/creditcard.png",
    color: "#6633CC",
    type: "form"
  }
]

export default function PaymentPage({ params }) {
  const router = useRouter()
  const { id: orderId } = use(params)
  const price = parseFloat(new URLSearchParams(window.location.search).get('price'))
  const [paymentMethod, setPaymentMethod] = useState('wechat')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [progress, setProgress] = useState(0)
  const [Order, setOrder] = useState(null)

  const [qrUrl, setQrUrl] = useState("");
  useEffect(() => {
    setQrUrl(`https://payment.example.com/qr/${Date.now()}`);
  }, []);

  useEffect(() => {
    if (loading) {
      const timer = setInterval(() => {
        setProgress((prev) => Math.min(prev + 33, 100))
      }, 500)
      return () => clearInterval(timer)
    }
  }, [loading])

  const handlePayment = async () => {
    try {
      setLoading(true)
      const reqbody = {
        order_id:parseInt(orderId),
        status:1
    }
      await apiClient.put(`/order/status`, reqbody)
      
      setSuccess(true)
      setTimeout(() => router.push(`/orders/${orderId}`), 2000)
    } catch (error) {
      console.error('Payment failed:', error)
    } finally {
      setLoading(false)
      setProgress(0)
    }
  }

  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-4">
      <div className="w-full max-w-4xl bg-white rounded-2xl shadow-xl overflow-hidden">
        {/* 进度条 */}
        <div className="h-1 bg-gray-200">
          <div 
            className="h-full transition-all duration-500" 
            style={{ 
              width: `${progress}%`,
              backgroundColor: paymentMethods.find(m => m.id === paymentMethod)?.color 
            }}
          />
        </div>

        <div className="grid md:grid-cols-2 gap-8 p-8">
          {/* 左侧支付信息 */}
          <div className="space-y-6">
            <div className="flex items-center gap-4">
              <Shield className="w-8 h-8 text-green-600" />
              <div>
                <h1 className="text-2xl font-bold">Secure Payment</h1>
                <p className="text-gray-500">Order: {orderId}</p>
              </div>
            </div>

            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="text-center">
                <p className="text-gray-600">Amount Due</p>
                <p className="text-4xl font-bold text-primary mt-2">
                  ￥{price.toFixed(2)}
                </p>
              </div>
            </div>

            <div className="space-y-4">
              <Label className="block text-lg font-medium">Payment Method</Label>
              <RadioGroup 
                value={paymentMethod} 
                onValueChange={setPaymentMethod}
                className="space-y-4"
              >
                {paymentMethods.map((method) => (
                  <div key={method.id}>
                    <RadioGroupItem 
                      value={method.id} 
                      id={method.id} 
                      className="sr-only"
                    />
                    <Label
                      htmlFor={method.id}
                      className={`flex items-center p-4 border-2 rounded-lg cursor-pointer transition-all
                        ${paymentMethod === method.id ? 'border-primary shadow-md' : 'border-gray-200'}`}
                    >
                      <img 
                        src={method.icon} 
                        alt={method.name}
                        className="w-10 h-10 mr-4"
                      />
                      <div className="flex-1">
                        <div className="font-medium">{method.name}</div>
                        {method.type === 'qrcode'}
                      </div>
                      {paymentMethod === method.id && (
                        <CheckCircle2 className="w-6 h-6 text-green-600" />
                      )}
                    </Label>
                  </div>
                ))}
              </RadioGroup>
            </div>
          </div>

          <div className="p-8 flex flex-col items-center justify-center bg-gray-50">
            {paymentMethods.find(m => m.id === paymentMethod)?.type === 'qrcode' ? (
              <div className="text-center space-y-6 w-full">
                <div className="inline-block p-4 bg-white rounded-lg shadow-sm border border-gray-200 mb-4">
                {qrUrl && (
                  <QRCode 
                    value={qrUrl} 
                    size={180}
                    level="H"
                    bgColor="#ffffff"
                    fgColor={paymentMethods.find(m => m.id === paymentMethod)?.color}
                  />
                )}
                </div>
                <div className="space-y-2 mb-6">
                  <p className="font-medium">
                    Scan with {paymentMethods.find(m => m.id === paymentMethod)?.name}
                  </p>
                  <p className="text-sm text-gray-500">
                    Valid for 15 minutes
                  </p>
                </div>
                <Button 
                  className="w-full max-w-xs"
                  onClick={handlePayment}
                  disabled={loading || success}
                  style={{ backgroundColor: paymentMethods.find(m => m.id === paymentMethod)?.color }}
                >
                  {loading ? (
                    <Loader2 className="w-4 h-4 animate-spin mr-2" />
                  ) : null}
                  {loading ? 'Processing...' : success ? 'Payment Confirmed' : 'Confirm Payment'}
                </Button>
              </div>
            ) : (
              <div className="w-full max-w-xs space-y-4">
                <div>
                  <Label>Card Number</Label>
                  <Input placeholder="4242 4242 4242 4242" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label>Expiry</Label>
                    <Input placeholder="MM/YY" />
                  </div>
                  <div>
                    <Label>CVV</Label>
                    <Input placeholder="123" type="password" />
                  </div>
                </div>
                <Button 
                  className="w-full mt-2"
                  onClick={handlePayment}
                  disabled={loading || success}
                >
                  {loading ? (
                    <Loader2 className="w-4 h-4 animate-spin mr-2" />
                  ) : null}
                  {loading ? 'Processing...' : success ? 'Payment Successful' : 'Pay Now'}
                </Button>
              </div>
            )}

            {success && (
              <div className="mt-6 text-center animate-in fade-in">
                <CheckCircle2 className="w-12 h-12 text-green-600 mx-auto mb-3" />
                <h3 className="text-lg font-semibold mb-1">
                  {paymentMethod === 'card' ? 'Payment Successful' : 'Payment Confirmed'}
                </h3>
                <p className="text-gray-500 text-sm">Redirecting to order details...</p>
              </div>
            )}

            <div className="mt-6 text-center text-xs text-gray-400">
              <div className="flex items-center justify-center gap-1">
                <Shield className="w-3 h-3" />
                Secure payment · Encrypted transaction
              </div>
            </div>
          </div>
        </div>

        {/* 安全提示
        <div className="border-t p-4 bg-gray-50">
          <div className="flex items-center justify-center gap-2 text-sm text-gray-500">
            <Shield className="w-4 h-4" />
            安全支付 · 平台担保交易 · 银行级数据加密
          </div>
        </div> */}
      </div>
    </div>
  )
}