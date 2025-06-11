'use client'

import { useState, useEffect } from 'react'
import apiClient from "@/lib/apiClient"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

export default function Profile() {
  const [user, setUser] = useState({
    name: '',
    email: '',
    password: '',
    is_seller: false,
    user_id: '',
  })

  const [addresses, setAddresses] = useState([])
  const [showPassword, setShowPassword] = useState(false)

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const res = await apiClient.get('/user'); // 调用后端 API 获取用户数据
        const { name, email, password, addresses } = res.data;

        setUser({
          name: name || '',
          email: email || '',
          password: password || '', // 密码通常不会从后端返回
          is_seller: localStorage.getItem('is_seller') === 'true' || false,
          user_id: localStorage.getItem('userId') || '',
        });

        setAddresses(addresses || []);
      } catch (error) {
        console.error('Error fetching user data:', error);
        alert('Failed to load user data');
      }
    };

    fetchUserData();
  }, []);

  const handleAddAddress = () => {
    setAddresses([...addresses, { address_name: '', address_phone: '', address: '' }])
  }
  
  const handleRemoveAddress = (index) => {
    const newAddresses = addresses.filter((_, i) => i !== index)
    setAddresses(newAddresses)
  }
  
  const handleAddressChange = (index, field, value) => {
    const newAddresses = addresses.map((address, i) => 
      i === index ? { ...address, [field]: value } : address
    )
    setAddresses(newAddresses)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      const requestBody = {
        name: user.name,
        email: user.email,
        password_hashed: user.password,
        addresses: addresses.map(addr => ({
          address_name: addr.address_name,
          address_phone: addr.address_phone,
          address: addr.address
        }))
      }

      const res = await apiClient.put('/user/update', requestBody)
      // localStorage.setItem('username', user.name)
      // localStorage.setItem('email', user.email)
      // localStorage.setItem('password', user.password)
      // localStorage.setItem('addresses', JSON.stringify(addresses))  

      alert('Profile updated successfully!')
    } catch (error) {
      console.error('Error updating profile:', error)
      alert('Error updating profile')
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">User Profile</h1>
      <Card>
        <CardHeader>
          <CardTitle>Edit Information</CardTitle>
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="grid w-full items-center gap-4">
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="name">Name</Label>
                <Input 
                  id="name" 
                  type="name"
                  value={user.name}
                  onChange={(e) => setUser({ ...user, name: e.target.value })}
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="email">Email</Label>
                <Input 
                  id="email" 
                  type="email"
                  value={user.email}
                  onChange={(e) => setUser({ ...user, email: e.target.value })}
                />
              </div>
              {user.is_seller && 
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="is_seller">SellerID</Label>
                <div>S
         S         <p>{user.user_id}</p>
                </div>
              </div>}
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="password">Password</Label>
                <Input 
                  id="password" 
                  type="password"
                  value={user.password}
                  onChange={(e) => setUser({ ...user, password: e.target.value })}
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label>Addresses</Label>
                {addresses.map((address, index) => (
                  <div key={index} className="p-4 border rounded-md space-y-2">
                    <div className="flex justify-between items-center">
                      <h3 className="text-lg font-semibold">Address {index + 1}</h3>
                      <Button type="button" variant="destructive" size="sm" onClick={() => handleRemoveAddress(index)}>Remove</Button>
                    </div>
                    <div className="flex flex-col space-y-1.5">
                      <Label htmlFor={`address_name-${index}`}>Name</Label>
                      <Input 
                        id={`address_name-${index}`}
                        value={address.address_name}
                        onChange={(e) => handleAddressChange(index, 'address_name', e.target.value)}
                      />
                    </div>
                    <div className="flex flex-col space-y-1.5">
                      <Label htmlFor={`address_phone-${index}`}>Phone</Label>
                      <Input 
                        id={`address_phone-${index}`}
                        value={address.address_phone}
                        onChange={(e) => handleAddressChange(index, 'address_phone', e.target.value)}
                      />
                    </div>
                    <div className="flex flex-col space-y-1.5">
                      <Label htmlFor={`address-${index}`}>Address</Label>
                      <Input 
                        id={`address-${index}`}
                        value={address.address}
                        onChange={(e) => handleAddressChange(index, 'address', e.target.value)}
                      />
                    </div>
                  </div>
                ))}
                <Button type="button" onClick={handleAddAddress}>Add Address</Button>
              </div>
            </div>
            <CardFooter className="mt-6">
              <Button type="submit">Update Profile</Button>
            </CardFooter>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}