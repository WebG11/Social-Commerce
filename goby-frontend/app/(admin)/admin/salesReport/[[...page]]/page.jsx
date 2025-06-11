"use client"

import { useState, useEffect } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import apiClient from "@/lib/apiClient";

export default function SalesReport() {
  const [salesData, setSalesData] = useState([]);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [topProducts, setTopProducts] = useState([]);
  const [reportData, setReportData] = useState(null);
  

  const fetchSalesReport = async () => {
    if (!startDate || !endDate) {
      alert("Please select both start and end dates.");
      return;
    }

    try {
      // 确保日期格式为YYYY-MM-DD
      const formattedStartDate = new Date(startDate).toISOString().split('T')[0];
      const formattedEndDate = new Date(endDate).toISOString().split('T')[0];
      
      const queryParams = `start_date=${formattedStartDate}&end_date=${formattedEndDate}`;
      const res = await apiClient.get(`/admin/order/report?${queryParams}`);

      // 使用 Object.entries 将 daily_revenue 对象转换为数组
      const dailyData = Object.entries(res.data.daily_revenue).map(([date, revenue]) => ({
        date,
        sales: revenue
      }));

      setSalesData(dailyData);
      setReportData(res.data);
      setTopProducts(res.data.top_products);
    } catch (error) {
      console.error('Error fetching sales report:', error);
    }
  };

  useEffect(() => {
    if (startDate && endDate) {
      fetchSalesReport();
    }
  }, [startDate, endDate]);

  return (
    <div>
      <div className="p-6 bg-white rounded-lg shadow-md">
        <h2 className="text-2xl font-bold mb-4 text-center">Sales Report</h2>
        <div className="flex justify-center mb-4">
          <input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className="p-2 border rounded mr-2"
          />
          <input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className="p-2 border rounded mr-2"
          />
          <button
            onClick={fetchSalesReport}
            className="p-2 bg-blue-500 text-white rounded"
          >
            Fetch Report
          </button>
        </div>
        {salesData.length > 0 && (
          <LineChart
            width={1200}
            height={400}
            data={salesData}
            margin={{
              top: 20, right: 50, left: 20, bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" interval={0} />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="sales" stroke="#8884d8" activeDot={{ r: 8 }} />
          </LineChart>
        )}
        {reportData && (
          <div className="mt-8 p-6 bg-white rounded-lg shadow-md">
            <h2 className="text-xl font-bold mb-4 text-center">Sales Summary</h2>
            <p>Average Order Amount: ￥{reportData.average_order_amt}</p>
            <p>Order Count: {reportData.order_count}</p>
            <p>Total Product Count: {reportData.total_product_count}</p>
            <p>Total Revenue: ￥{reportData.total_revenue}</p>
            <h3 className="mt-4 font-bold">Best Selling Product:</h3>
            {(() => {
              // 将对象转换为数组并获取最后一个（销量最好的）
              const topProductsArray = Object.entries(reportData.top_products);
              if (topProductsArray.length > 0) {
                const [productName, count] = topProductsArray[topProductsArray.length - 1];
                return (
                  <div className="mt-2 p-4 border rounded-lg">
                    <p className="text-lg font-semibold">{productName}</p>
                    <p>Sold: {count} units</p>
                  </div>
                );
              }
              return <p>No products data available</p>;
            })()}
          </div>
        )}
      </div>
    </div>
  );
}