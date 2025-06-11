"use client"

import { Button } from "@/components/ui/button"
import { useState } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import apiClient from "@/lib/apiClient"
import ImageUploader from "@/components/imageUploader"
import { redirect } from 'next/navigation'
import AddressDropdown from "./addressDropdown"
import { useParams, useRouter } from 'next/navigation'

export default ({ selectedItems, total_price })=>{
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [name,setName] = useState("")
    const [price,setPrice] = useState(0.00)
    const [stock,setStock] = useState(0)
    const [description,setDescription] = useState("");
    const [imageUrl, setImageUrl] = useState(null)

    const router = useRouter()
    const handleSubmit = async () => {
      const ids = selectedItems.map(item => item.id);
      const res = await apiClient.post('/order', { itemIDs: ids, total_price: total_price });
      console.log(res);
      if(res.status === 200){
        alert('Order submitted successfully');
        const order = res.data.order;
        router.push(`/orders/${order.id}`)
      }else{
        alert('订单中的商品必须来自同一个卖家');
        alert(res.data.message);
      }
    }
    return (
        <>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogTrigger asChild>
                    <Button className="mb-4">Create New Product</Button>
                </DialogTrigger>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Select Shipping Address</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        {/* <Button disabled={selectedItems.length === 0 } onClick={handleOpen}>提交订单</Button> */}
                        <AddressDropdown/>
                    </div>
                </DialogContent>
            </Dialog>
        </>
       
    )
}