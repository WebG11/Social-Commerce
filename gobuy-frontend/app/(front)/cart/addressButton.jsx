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

export default ({ selectedItems })=>{
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [name,setName] = useState("")
    const [price,setPrice] = useState(0.00)
    const [stock,setStock] = useState(0)
    const [description,setDescription] = useState("");
    const [imageUrl, setImageUrl] = useState(null)

    const router = useRouter()
    const handleSubmit = async () => {
      const ids = selectedItems.map(item => item.id);
      const res = await apiClient.post('/order', { itemIDs: ids });
      alert('提交成功');
      const order = res.data.order;
      router.push(`/orders/${order.id}`)
    }
    return (
        <>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogTrigger asChild>
                    <Button className="mb-4">新建商品</Button>
                </DialogTrigger>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>选择收货地址</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        <AddressDropdown/>
      <Button disabled={selectedItems.length === 0 } onClick={handleCreate}>提交订单</Button>

                    </div>
                </DialogContent>
            </Dialog>
        </>
       
    )
}