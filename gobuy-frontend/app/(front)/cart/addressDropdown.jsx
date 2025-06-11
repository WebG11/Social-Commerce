"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { PlusCircle, MapPin } from "lucide-react"
import { Label } from "@/components/ui/label"
import apiClient from "@/lib/apiClient"
import { Input } from "@/components/ui/input"

function AddressList({ addresses, selectedAddress, onSelect, onEdit, onDelete }) {
  if (addresses.length === 0) {
    return <p className="text-center text-gray-500">暂无收货地址</p>
  }

  return (
    <ul className="space-y-2">
      {addresses.map((address) => (
        <li
          key={address.id}
          className={`flex items-center justify-between p-2 border rounded cursor-pointer ${
            selectedAddress?.id === address.id ? "bg-primary/10" : ""
          }`}
          onClick={() => onSelect(address)}
        >
          <div>
            <p className="font-medium">{address.name}</p>
            <p className="text-sm text-gray-600">{address.phone}</p>
            <p className="text-sm text-gray-600">{address.address}</p>
          </div>
          <div onClick={(e) => e.stopPropagation()}>
            <Button variant="ghost" size="icon" onClick={() => onEdit(address)}>
              <Pencil className="h-4 w-4" />
            </Button>
            <Button variant="ghost" size="icon" onClick={() => onDelete(address.id)}>
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </li>
      ))}
    </ul>
  )
}

 function AddressForm({ address, onSubmit, onCancel }) {
  const [name, setName] = useState(address?.name || "")
  const [phone, setPhone] = useState(address?.phone || "")
  const [addressText, setAddressText] = useState(address?.address || "")

  const handleSubmit = (e) => {
    e.preventDefault()
    onSubmit({ name, phone, address: addressText })
    setName("")
    setPhone("")
    setAddressText("")
  }

  return (
    <form className="space-y-4">
      <div>
        <Label htmlFor="name">收货人姓名</Label>
        <Input id="name" value={name} onChange={(e) => setName(e.target.value)} required />
      </div>
      <div>
        <Label htmlFor="phone">联系电话</Label>
        <Input id="phone" value={phone} onChange={(e) => setPhone(e.target.value)} required />
      </div>
      <div>
        <Label htmlFor="address">详细地址</Label>
        <Input id="address" value={addressText} onChange={(e) => setAddressText(e.target.value)} required />
      </div>
      <div className="flex justify-end space-x-2">
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            取消
          </Button>
        )}
        <Button onClick={handleSubmit} type="button">{address ? "更新" : "添加"}</Button>
      </div>
    </form>
  )
}

export default function AddressDropdown() {
  const [addresses, setAddresses] = useState([])
  const [selectedAddress, setSelectedAddress] = useState(null)
  const [editingAddress, setEditingAddress] = useState(null)
  const [isAddingNew, setIsAddingNew] = useState(false)

  const addAddress = (newAddress) => {
    const address = { ...newAddress, id: Date.now().toString() }
    setAddresses([...addresses, address])
    setSelectedAddress(address)
    setIsAddingNew(false)
  }

  const updateAddress = (updatedAddress) => {
    setAddresses(addresses.map((addr) => (addr.id === updatedAddress.id ? updatedAddress : addr)))
    setSelectedAddress(updatedAddress)
    setEditingAddress(null)
  }

  const deleteAddress = (id) => {
    setAddresses(addresses.filter((addr) => addr.id !== id))
    if (selectedAddress?.id === id) {
      setSelectedAddress(null)
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" className="w-full justify-start text-left font-normal">
          <MapPin className="mr-2 h-4 w-4" />
          <span>{selectedAddress ? `${selectedAddress.name}, ${selectedAddress.phone}` : "选择收货地址"}</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>收货地址</DialogTitle>
        </DialogHeader>
        <div className="mt-4 space-y-4">
          <AddressList
            addresses={addresses}
            selectedAddress={selectedAddress}
            onSelect={setSelectedAddress}
            onEdit={setEditingAddress}
            onDelete={deleteAddress}
          />
          {!isAddingNew && !editingAddress && (
            <Button onClick={() => setIsAddingNew(true)} className="w-full">
              <PlusCircle className="mr-2 h-4 w-4" />
              新建地址
            </Button>
          )}
          {isAddingNew && <AddressForm onSubmit={addAddress} onCancel={() => setIsAddingNew(false)} />}
          {editingAddress && (
            <AddressForm address={editingAddress} onSubmit={updateAddress} onCancel={() => setEditingAddress(null)} />
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

