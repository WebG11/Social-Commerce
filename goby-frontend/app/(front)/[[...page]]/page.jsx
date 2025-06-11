'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams, useParams  } from 'next/navigation'
import Link from 'next/link'
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import SearchInput from '@/components/productSearchInput';
import Pagination from '@/components/pagination';
import AgentSearch from '@/components/agentSearch';
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select"
import apiClient from '@/lib/apiClient';

export default function Page() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const params = useParams()

  const [products, setProducts] = useState([])
  const [promotions, setPromotions] = useState([])
  const [totalPage, setTotalPage] = useState(1)

  // 获取参数
  const category = searchParams.get('category') || ""
  const search = searchParams.get('search') || ""
  const currentPage = Number(params.page?.[0] || 1)
  const productsPerPage = 10
  const minPrice = searchParams.get('min_price') || ''
  const maxPrice = searchParams.get('max_price') || ''
  const brand = searchParams.get('brand') || ''
  
  // 数据获取函数
  const fetchProducts = async () => {
    try {
      let page = search ? 1 : currentPage;
      if (isNaN(page)) {
        page = 1;
      }
      const queryParams = new URLSearchParams({
        page: page,
        size: productsPerPage,
        query: search,
        category: category,
        min_price: minPrice,
        max_price: maxPrice,
        brand: brand,
      }).toString();
      const res = await apiClient.get(`/product/search?${queryParams}`);
      setProducts(res.data.products)
      setTotalPage(Math.ceil(res.data.total_count / productsPerPage))
    } catch (error) {
      console.error("Failed to fetch data:", error)
    }
  }

  const fetchPromotions = async () => {
    try {
      const res = await apiClient.get('/product/promotions');
      setPromotions(res.data.promotions);
    } catch (error) {
      console.error('Error fetching promotions:', error);
    }
  }

  // 监听参数变化
  useEffect(() => {
    fetchProducts()
    fetchPromotions()
  }, [currentPage, search, category, minPrice, maxPrice, brand])

  // 分类选择处理
  const handleCategoryChange = (value) => {
    const newParams = new URLSearchParams(searchParams)
    if (value === "all") {
      newParams.delete('category') // 清除 category 参数
    } else {
      newParams.set('category', value) // 设置新的 category 参数
    }
    router.replace(`?${newParams.toString()}`) // 无刷新更新URL
  }

  // 价格和品牌筛选处理
  const handleMinPriceChange = (e) => {
    const newParams = new URLSearchParams(searchParams)
    if (e.target.value === '') {
      newParams.delete('min_price')
    } else {
      newParams.set('min_price', e.target.value)
    }
    router.replace(`?${newParams.toString()}`)
  }
  const handleMaxPriceChange = (e) => {
    const newParams = new URLSearchParams(searchParams)
    if (e.target.value === '') {
      newParams.delete('max_price')
    } else {
      newParams.set('max_price', e.target.value)
    }
    router.replace(`?${newParams.toString()}`)
  }
  const handleBrandChange = (e) => {
    const newParams = new URLSearchParams(searchParams)
    if (e.target.value === '') {
      newParams.delete('brand')
    } else {
      newParams.set('brand', e.target.value)
    }
    router.replace(`?${newParams.toString()}`)
  }

  const categories = [
    { id: 'Electronics', name: 'Electronics' }, // 电子产品
    { id: 'Home & Furniture', name: 'Home & Furniture' },  // 家居与家具
    { id: 'Apparel & Accessories', name: 'Apparel & Accessories' },  // 服装与配饰
    { id: 'Food & Beverage', name: 'Food & Beverage' },  // 食品与饮料
    { id: 'Beauty & Cosmetics', name: 'Beauty & Cosmetics' },  // 美妆与护肤品
    { id: 'Toys & Games', name: 'Toys & Games' },  // 玩具与游戏
    { id: 'Books & Stationery', name: 'Books & Stationery' },  // 书籍与文具
    { id: 'Sports & Outdoors', name: 'Sports & Outdoors' },  // 运动与户外
  ];

  return (  
    <>
      <div>
        <div className='flex gap-2 mb-2 flex-wrap'>
          {/* 新增分类选择器 */}
          <div className="w-48">
            <Select
              value={category}
              onValueChange={handleCategoryChange}
            >
              <SelectTrigger>
                <SelectValue placeholder="ALL" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">ALL</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={cat.id}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {/* 新增价格区间筛选 */}
          <input
            type="number"
            min="0"
            placeholder="Min Price"
            value={minPrice}
            onChange={handleMinPriceChange}
            className="w-24 p-2 border rounded text-sm"
          />
          <span className="self-center text-sm">-</span>
          <input
            type="number"
            min="0"
            placeholder="Max Price"
            value={maxPrice}
            onChange={handleMaxPriceChange}
            className="w-24 p-2 border rounded text-sm"
          />
          {/* 新增品牌筛选 */}
          <input
            type="text"
            placeholder="Brand"
            value={brand}
            onChange={handleBrandChange}
            className="w-32 p-2 border rounded text-sm"
          />
          <SearchInput currentPage={currentPage} initialSearch={search} className="text-sm"/>
          {/* <AgentSearch className="text-sm"/> */}
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-5 gap-6">
          {products.map((product) => {
            const productPromotion = promotions.find(promo => parseInt(promo.product_id) === product.id);
            const discountedPrice = productPromotion ? product.price * (1 - productPromotion.discount_rate / 100) : product.price;

            return (
              <Link href={`/product/${product.id}`}  key={product.id}>
                <Card>
                  <CardContent className="p-0">
                    <img
                      src={product.image || "/placeholder.svg"}
                      alt={product.name}
                      className="w-full h-full object-cover aspect-square"
                    />
                  </CardContent>
                  <CardFooter className="flex flex-col items-start">
                    <h3 className="font-semibold line-clamp-1">{product.name}</h3>
                    <p className="text-sm text-muted-foreground line-clamp-2">{product.description}</p>
                    <p className="text-xl font-bold mt-1 text-orange-500">
                      {productPromotion ? (
                        <>
                          <span className="line-through text-gray-500">￥{product.price.toFixed(2)}</span>
                          <span className="text-red-500 ml-2">￥{discountedPrice.toFixed(2)}</span>
                        </>
                      ) : (
                        <>￥{product.price.toFixed(2)}</>
                      )}
                    </p>
                    <p className="text-sm font-medium mt-1">Stock: {product.stock} items</p>
                  </CardFooter>
                </Card>
              </Link>
            );
          })}
        </div>
        <Pagination 
          pageSize={productsPerPage}
          currentPage={currentPage}
          totalPage={totalPage}
          hrefPrefix={``} 
          hrefParam={{ 
            search: search,
            category: category,
          }}
          />
      </div>
    </>
  )
}