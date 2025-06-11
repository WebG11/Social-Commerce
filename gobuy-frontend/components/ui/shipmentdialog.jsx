"use client"
import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export function ShipmentDialog({ order, open, onOpenChange, onConfirm }) {
  const [trackingNumber, setTrackingNumber] = useState('')

  const handleSubmit = async () => {
    if (!trackingNumber) return
    await onConfirm(order.id, trackingNumber)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Shipping Information</DialogTitle>
          <DialogDescription>
            Please fill in the shipping information to complete the shipment process
          </DialogDescription>
        </DialogHeader>
        {order && (
        <div className="space-y-4">
          <div>
            <Label>Order Number</Label>
            <Input value={order.number} disabled />
          </div>
          <div>
            <Label>Tracking Number</Label>
            <Input 
              value={trackingNumber}
              onChange={(e) => setTrackingNumber(e.target.value)}
              placeholder="Enter tracking number"
              required
            />
          </div>
        </div>
        )}
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
          Cancel
          </Button>
          <Button onClick={handleSubmit}>Confirm Shipment</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}