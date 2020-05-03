import httpClient from '../../../../utils/httpClient'
import {v4 as uuidv4} from 'uuid'

export const getOrders = products => products.reduce((acc, product) => ({...acc, [product.id]: {...product, quantity: (acc[product.id] && ++acc[product.id].quantity || 1)}}),{})


const submitOrder = async (order, groupID) => {
  const url = new URL('/api/v1/orders', process.env.REACT_APP_PRODUCTS_API)
  const {status,data} = await httpClient({
    method: 'POST',
    url: url.toString(),
    data: JSON.stringify({
      uid: groupID,
      productId: order.id,
      quantity: order.quantity,
      userId: order.userId,
    })
  })
  
  if(status !== 200){
    throw new Error(data.message)
  }
}

const submitOrders = async (orders, userProfile) => {
  const groupID = uuidv4()
  const mappingFn = order => submitOrder({...order, userId: userProfile.id},groupID)
  return Promise.all(orders.map(mappingFn))
}


export default submitOrders