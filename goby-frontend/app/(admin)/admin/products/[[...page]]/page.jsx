'use client'

import { useState, useEffect, use } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import Pagination from '@/components/pagination'
import apiClient from "@/lib/apiClient"
import NewButton from "./newButton"
import DeleteButton from "./deleteButton"
import UpdateButton from "./updateButton"

export default function ProductManagement({ searchParams, params }) {
  const productsPerPage = 10
  const resolvedParams = use(params) // 解包 params
  const currentPage = Number(resolvedParams?.page || 1)

  const [products, setProducts] = useState([])
  const [totalPage, setTotalPage] = useState(0)

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const res = await apiClient.get("/admin/product/list", { page: currentPage, size: productsPerPage })
        const { products, total_count } = res.data
        setProducts(products)
        setTotalPage(Math.ceil(total_count / productsPerPage))
      } catch (error) {
        console.error("Failed to fetch products:", error)
      }
    }

    fetchProducts()
  }, [currentPage])

  const handleDeleteProduct = (productId) => {
    setProducts(products.filter((p) => p.id !== productId))
  }

  const handleUpdateProduct = (updatedProduct) => {
    setProducts((prevProducts) =>
      prevProducts.map((product) =>
        product.id === updatedProduct.id ? updatedProduct : product
      )
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Product Management</h2>
      <div className="flex gap-2">
        <NewButton />
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Image</TableHead>
            <TableHead>Price</TableHead>
            <TableHead>Description</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {products.map((product) => (
            <TableRow key={product.id}>
              <TableCell>{product.id}</TableCell>
              <TableCell>{product.name}</TableCell>
              <TableCell>
                <img src={product.image} className="h-10 w-10" />
              </TableCell>
              <TableCell>￥{product.price.toFixed(2)}</TableCell>
              <TableCell>{product.description}</TableCell>
              <TableCell>
                <div className="flex flex-col gap-1">
                <UpdateButton 
                  product={product}
                  // onUpdate={(updatedProduct) => {
                  //   // 可选：本地更新逻辑
                  //   setProducts(products.map(p => 
                  //     p.id === updatedProduct.id ? updatedProduct : p
                  //   ))
                  // }}
                  onDelete={() => handleUpdateProduct(product)} 
                />
                  <DeleteButton productId={product.id} onDelete={() => handleDeleteProduct(product.id)} />
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Pagination hrefPrefix={`/admin/products`} currentPage={currentPage} pageSize={productsPerPage} totalPage={totalPage} />
    </div>
  )
}