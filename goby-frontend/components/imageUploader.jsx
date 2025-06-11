"use client"

import { useState, useRef } from "react"
import { Button } from "@/components/ui/button"
import Image from "next/image"
import uploadImage from "../lib/uploadImage"

export default  ({setImageUrl})=> {
  
  const [isUploading, setIsUploading] = useState(false)
  const fileInputRef = useRef(null)

  const handleUpload = async (event) => {
    const file = event.target.files?.[0]
    if (!file) return

    // 判断文件大小是否小于等于 5MB
    const maxSize = 5 * 1024 * 1024; // 5MB 转换为字节
    if (file.size > maxSize) {
      alert('File size cannot exceed 5MB');
      return;
    }

    setIsUploading(true)

    const formData = new FormData()
    formData.append("file", file)

    try {
      const url = await uploadImage(formData)
      setImageUrl(url)
    } catch (error) {
      console.error("Upload failed:", error)
      // Here you might want to show an error message to the user
    } finally {
      setIsUploading(false)
    }
  }

  const handleClick = () => {
    fileInputRef.current?.click()
  }

  return (
    <div className="flex space-y-4">
      <Button onClick={handleClick} disabled={isUploading}>
        {isUploading ? "Uploading..." : "Upload Image"}
      </Button>
      <input type="file" ref={fileInputRef} onChange={handleUpload} accept="image/*" className="hidden" />
    </div>
  )
}

