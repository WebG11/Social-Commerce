'use client';

import { useState, useEffect } from "react";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { format } from 'date-fns';
import apiClient from "@/lib/apiClient";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";

export default function PromotionsPage() {
  const [promotions, setPromotions] = useState([]);
  const [products, setProducts] = useState([]);
  const [isDiscountDialogOpen, setDiscountDialogOpen] = useState(false);
  const [isCouponDialogOpen, setCouponDialogOpen] = useState(false);
  const [formData, setFormData] = useState({
    product_id: '',
    discount_rate: '',
    coupon_code: '',
    coupon_type: 'percentage',
    coupon_value: '',
    min_purchase: '',
    usage_limit: '',
    start_date: '',
    end_date: ''
  });

  useEffect(() => {
    fetchPromotions();
  }, []);

  const fetchProducts = async () => {
    try {
      const res = await apiClient.get('/admin/product/listall');
      setProducts(res.data.products);
    } catch (error) {
      console.error('Error fetching products:', error);
    }
  };

  const fetchPromotions = async () => {
    try {
      const res = await apiClient.get('/product/promotions');
      setPromotions(res.data.promotions);
    } catch (error) {
      console.error('Error fetching promotions:', error);
    }
  };
  const handleDeletePromotion = async (promotionId) => {
    try {
      await apiClient.delete(`/product/promotions/${promotionId}`);
      alert('Promotion deleted successfully');
      fetchPromotions(); // 重新获取促销活动列表以更新UI
    } catch (error) {
      console.error('Error deleting promotion:', error);
      alert('Failed to delete promotion');
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmitDiscount = async (e) => {
    e.preventDefault();
    try {
      await apiClient.post('/product/promotions', {
        type: 'discount',
        ...formData
      });
      setDiscountDialogOpen(false);
      setFormData({
        product_id: '',
        discount_rate: '',
        start_date: '',
        end_date: ''
      });
      fetchPromotions();
    } catch (error) {
      console.error('Error creating discount:', error);
    }
  };

  const handleSubmitCoupon = async (e) => {
    e.preventDefault();
    try {
      await apiClient.post('/product/promotions', {
        type: 'coupon',
        ...formData
      });
      setCouponDialogOpen(false);
      setFormData({
        coupon_code: '',
        coupon_type: 'percentage',
        coupon_value: '',
        min_purchase: '',
        usage_limit: '',
        start_date: '',
        end_date: ''
      });
      fetchPromotions();
    } catch (error) {
      console.error('Error creating coupon:', error);
    }
  };

  const openDiscountDialog = () => {
    fetchProducts(); // 确保在打开对话框之前调用
    setDiscountDialogOpen(true);
  }

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Promotion Management</h2>
      <div className="flex gap-2 mb-4">
        <Dialog open={isDiscountDialogOpen} onOpenChange={setDiscountDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openDiscountDialog}>Create Product Discount</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Product Discount</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmitDiscount} className="space-y-2">
              <div>
                <label className="block">Product</label>
                <select
                  name="product_id"
                  value={formData.product_id}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                >
                  <option value="">Select Product</option>
                  {products.map(product => (
                    <option key={product.id} value={product.id}>
                      {product.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block">Discount Rate (%)</label>
                <input
                  type="number"
                  name="discount_rate"
                  value={formData.discount_rate}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  min="1"
                  max="100"
                  required
                />
              </div>
              <div>
                <label className="block">Start Date</label>
                <input
                  type="date"
                  name="start_date"
                  value={formData.start_date}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                />
              </div>
              <div>
                <label className="block">End Date</label>
                <input
                  type="date"
                  name="end_date"
                  value={formData.end_date}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                />
              </div>
              <div className="flex gap-2">
                <Button type="submit">Create</Button>
                <Button type="button" variant="outline" onClick={() => setDiscountDialogOpen(false)}>Cancel</Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>

        <Dialog open={isCouponDialogOpen} onOpenChange={setCouponDialogOpen}>
          <DialogTrigger asChild>
            <Button>Create Coupon</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Coupon</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmitCoupon} className="space-y-2">
              <div>
                <label className="block">Coupon Code</label>
                <input
                  type="text"
                  name="coupon_code"
                  value={formData.coupon_code}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                />
              </div>
              <div>
                <label className="block">Coupon Type</label>
                <select
                  name="coupon_type"
                  value={formData.coupon_type}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                >
                  <option value="percentage">Percentage Discount</option>
                  <option value="fixed">Fixed Amount</option>
                </select>
              </div>
              {formData.coupon_type === 'fixed' ? (
                <>
                  <div>
                    <label className="block">Coupon Value</label>
                    <input
                      type="number"
                      name="coupon_value"
                      value={formData.coupon_value}
                      onChange={handleInputChange}
                      className="w-full p-2 border rounded"
                      min="1"
                      required
                    />
                  </div>
                  <div>
                    <label className="block">Minimum Purchase Amount</label>
                    <input
                      type="number"
                      name="min_purchase"
                      value={formData.min_purchase}
                      onChange={handleInputChange}
                      className="w-full p-2 border rounded"
                      min="0"
                      required
                    />
                  </div>
                </>
              ) : (
                <>
                  <div>
                    <label className="block">Maximum Discount Amount</label>
                    <input
                      type="number"
                      name="usage_limit"
                      value={formData.usage_limit}
                      onChange={handleInputChange}
                      className="w-full p-2 border rounded"
                      min="1"
                      required
                    />
                  </div>
                  <div>
                    <label className="block">Discount Rate (%)</label>
                    <input
                      type="number"
                      name="coupon_value"
                      value={formData.coupon_value}
                      onChange={handleInputChange}
                      className="w-full p-2 border rounded"
                      min="1"
                      max="100"
                      required
                    />
                  </div>
                  <div>
                    <label className="block">Minimum Purchase Amount</label>
                    <input
                      type="number"
                      name="min_purchase"
                      value={formData.min_purchase}
                      onChange={handleInputChange}
                      className="w-full p-2 border rounded"
                      min="0"
                      required
                    />
                  </div>
                </>
              )}
              <div>
                <label className="block">Start Date</label>
                <input
                  type="date"
                  name="start_date"
                  value={formData.start_date}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                />
              </div>
              <div>
                <label className="block">End Date</label>
                <input
                  type="date"
                  name="end_date"
                  value={formData.end_date}
                  onChange={handleInputChange}
                  className="w-full p-2 border rounded"
                  required
                />
              </div>
              <div className="flex gap-2">
                <Button type="submit">Create</Button>
                <Button type="button" variant="outline" onClick={() => setCouponDialogOpen(false)}>Cancel</Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Type</TableHead>
            <TableHead>Details</TableHead>
            <TableHead>Validity</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Operation</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {promotions.map((promotion) => (
            <TableRow key={promotion.id}>
              <TableCell>{promotion.type === 'discount' ? 'Product Discount' : 'Coupon'}</TableCell>
              <TableCell>
                {promotion.type === 'discount' ? (
                  <div>
                    <div>Product ID: {promotion.product_id}</div>
                    <div>Discount Rate: {promotion.discount_rate}%</div>
                  </div>
                ) : (
                  <div>
                    <div>Code: {promotion.coupon_code}</div>
                    <div>Type: {promotion.coupon_type === 'percentage' ? 'Percentage Discount' : 'Fixed Amount'}</div>
                    <div>Value: {promotion.coupon_value}</div>
                    <div>Minimum Purchase: {promotion.min_purchase}</div>
                  </div>
                )}
              </TableCell>
              <TableCell>
                <div>Start Date: {format(new Date(promotion.start_date), 'yyyy-MM-dd')}</div>
                <div>End Date: {format(new Date(promotion.end_date), 'yyyy-MM-dd')}</div>
              </TableCell>
              <TableCell>
                <Badge variant={new Date(promotion.end_date) > new Date() ? 'success' : 'destructive'}>
                  {new Date(promotion.end_date) > new Date() ? 'Ongoing' : 'Expired'}
                </Badge>
              </TableCell>
              <TableCell>
                <Button variant="destructive" onClick={() => handleDeletePromotion(promotion.id)}>Delete</Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}