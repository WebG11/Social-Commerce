"use client"

import { Button } from "@/components/ui/button"

export default ({status})=>{
    return (
        <Button 
            variant={status === 'Active' ? 'destructive' : 'default'}
        >
            {status === 'Active' ? 'Ban' : 'Unban'}
        </Button>
    )
}