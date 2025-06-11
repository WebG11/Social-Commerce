'use server'
 
import apiClient from '@/lib/apiClient';
import { redirect } from 'next/navigation'
export async function handleSearch(currentPage,search) {
    const queryString = '?' + new URLSearchParams({search}).toString();
    redirect(`/${currentPage}${queryString}`) 
}

export async function handleAskAgent(input) {
    const res = await apiClient.post("/agent/ask",{"user_prompt":input});
    return {
        html: res.data.content
    }
}