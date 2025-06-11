import Image from "next/image"
import { notFound } from "next/navigation"
import apiClient from '@/lib/apiClient';
import AddButton from "./addButton";
import { FaStar } from 'react-icons/fa';

export default async function ProductPage({ params }) {
  const { id } = await params
  const res = await apiClient.get(`/product/${id}`)
  const product  = res.data.product;
  const reviewsRes = await apiClient.get(`/product/${id}/reviews`);
  const reviews = reviewsRes.data.reviews;

  // 获取所有促销活动
  const promotionsRes = await apiClient.get('/product/promotions');
  const promotions = promotionsRes.data.promotions;

  // 筛选出与当前商品相关的促销活动
  const productPromotion = promotions.find(promo => promo.product_id === id);

  // 计算折扣价
  const discountedPrice = productPromotion ? product.price * (1 - productPromotion.discount_rate / 100) : product.price;

  if (!product) {
    notFound()
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid md:grid-cols-2 gap-8">
        <div className="relative aspect-square">
          <Image
            src={product.image || null}
            alt={product.name}
            layout="fill"
            objectFit="cover"
            className="rounded-lg"
          />
        </div>
        <div className="flex flex-col justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2">{product.name}</h1>
            <p className="text-2xl font-semibold mb-4">
              {productPromotion ? (
                <>
                  <span className="line-through text-gray-500">￥{product.price.toFixed(2)}</span>
                  <span className="text-red-500 ml-2">￥{discountedPrice.toFixed(2)}</span>
                </>
              ) : (
                <>￥{product.price.toFixed(2)}</>
              )}
            </p>
            <p className="text-gray-600 mb-0">Stock: {product.stock} items</p>
            <p className="text-gray-600 mb-4">Seller: {product.seller_id}</p>
            <div className="prose max-w-none mb-6">
              <p>{product.description}</p>
            </div>
            <div className="mb-6">
              <h2 className="text-xl font-semibold mb-2">Reviews</h2>
              {reviews.length > 0 ? (
                reviews.map((review, index) => (
                  <div key={index} className="mb-4">
                    <div className="flex items-center">
                      {[...Array(5)].map((_, i) => (
                        <FaStar
                          key={i}
                          size={16}
                          color={i < review.rating ? "#ffc107" : "#e4e5e9"}
                        />
                      ))}
                    </div>
                    <p className="text-gray-600">{review.comment}</p>
                    <div className="text-sm text-gray-500">
                      <span>By {review.user_name}</span> | <span>{review.created_at}</span>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-gray-600">No reviews yet.</p>
              )}
            </div>
          </div>
          <AddButton productID={product.id} sellerID={product.seller_id}/>
        </div>
      </div>
    </div>
  )
}