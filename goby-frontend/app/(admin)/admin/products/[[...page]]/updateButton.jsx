"use client"

import { Button } from "@/components/ui/button"
import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import ImageUploader from "@/components/imageUploader"
import apiClient from "@/lib/apiClient"
import { useRouter } from "next/navigation"
import { toast } from "sonner"

export default ({ product, onUpdate }) => {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    id: product.id,
    name: product.name,
    price: product.price,
    stock: parseInt(product.stock, 10),
    description: product.description,
    image: product.image
  })
  const router = useRouter()

  const handleUpdate = async () => {
    try {
      setLoading(true)
      // 确保 stock 是整数形式
        const updatedFormData = {
            ...formData,
            price: parseFloat(formData.price),
            stock: parseInt(formData.stock, 10), // 将 stock 转换为整数
        }

      const res = await apiClient.put(`/admin/product/${product.id}`, updatedFormData)
      toast.success("Product updated successfully")
      onUpdate?.(res.data)
      router.refresh()
      setOpen(false)
    } catch (error) {
      toast.error("Update failed: " + (error.response?.data?.message || error.message))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button 
          className="px-2 py-1 text-sm h-auto"
          disabled={loading}
        >
          {loading ? "Saving..." : "Edit"}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Edit Product</DialogTitle>
        </DialogHeader>
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <div>
              <Label>Product Name</Label>
              <Input
                value={formData.name}
                onChange={e => setFormData({...formData, name: e.target.value})}
              />
            </div>
            <div>
              <Label>Price</Label>
              <Input
                type="number"
                step="0.01"
                value={formData.price}
                onChange={e => setFormData({...formData, price: e.target.value})}
              />
            </div>
            <div>
              <Label>Stock Quantity</Label>
              <Input
                type="number"
                value={formData.stock}
                onChange={e => setFormData({...formData, stock: e.target.value})}
              />
            </div>
          </div>
          <div className="space-y-4">
            <div>
              <Label>Description</Label>
              <Input
                value={formData.description}
                onChange={e => setFormData({...formData, description: e.target.value})}
              />
            </div>
            <div>
              <Label>Product Image</Label>
              <ImageUploader 
                initialImage={formData.image}
                onUpload={url => setFormData({...formData, image: url})}
              />
            </div>
          </div>
        </div>
        <Button 
          className="w-full mt-4"
          onClick={handleUpdate}
          disabled={loading}
        >
          {loading ? "Saving..." : "Save Changes"}
        </Button>
      </DialogContent>
    </Dialog>
  )
}