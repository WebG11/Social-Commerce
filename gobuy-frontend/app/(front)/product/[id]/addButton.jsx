'use client'

import { Button } from "@/components/ui/button"
import apiClient from "@/lib/apiClient"


export default ({productID, sellerID})=>{
    const handleAdd = async()=>{
        if (localStorage.getItem("is_seller") === "true" && parseInt(localStorage.getItem("userId")) === parseInt(sellerID)){
            alert("You are not allowed to buy your own product");
            return;
        }
        const res = await apiClient.post(`/cart/${productID}`);
        alert("Added to cart successfully");
    }
    return (
        <Button onClick={handleAdd} className="bg-orange-500 text-white py-2 px-4 w-52 rounded transition-colors hover:bg-orange-600" >
            Add to Cart
        </Button>
    )
}