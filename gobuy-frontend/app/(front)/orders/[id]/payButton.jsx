"use client"
import { Button } from "@/components/ui/button"
import { redirect } from "next/navigation"
import { useRouter } from "next/navigation"
import apiClient from "@/lib/apiClient"

export default ({ orderID, address, price }) => {
  const router = useRouter()

  const handlePayClick = async() => {
            const reqbody = {
            order_id:parseInt(orderID),
            address_name:address.address_name,
            address_phone:address.address_phone,
            address:address.address
        }
    const addrres = await apiClient.put(`/order/orderAddress`,reqbody)

    // 跳转到支付页面并传递订单ID
    router.push(`/payment/${orderID}?price=${price}`)
  }

  const handlePay = async()=>{
    const reqbody = {
      order_id:parseInt(orderID),
      address_name:address.address_name,
      address_phone:address.address_phone,
      address:address.address
    }
    const addrres = await apiClient.put(`/order/orderAddress`,reqbody)

    const res = await apiClient.get(`/payment/${orderID}`,{order_id:orderID})
    redirect(res.data.url);
  }

  return (
    <Button onClick={handlePay} className="bg-green-500 hover:bg-green-600">
      Pay
    </Button>
  )
}
// "use client"
// import { Button } from "@/components/ui/button"
// import apiClient from "@/lib/apiClient"
// import { redirect } from "next/navigation"

// export default ({orderID, address})=>{
//     const handlePay = async()=>{
//         const reqbody = {
//             order_id:parseInt(orderID),
//             address_name:address.address_name,
//             address_phone:address.address_phone,
//             address:address.address
//         }
//         const res = await apiClient.get(`/payment/${orderID}`,{order_id:orderID})
//         const addrres = await apiClient.put(`/order/orderAddress`,reqbody)
//         redirect(res.data.url);
//     }
//     return (
//         <Button onClick={handlePay}>Pay</Button>
//     )
// }