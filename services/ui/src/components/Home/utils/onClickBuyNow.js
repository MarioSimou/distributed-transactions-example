import submitOrders, {getOrders} from '../components/utils/submitOrders.js'

const onClickBuyNow = (products, userProfile, resetCart) => async () => {
  console.log(products)
  if(products.length === 0) {
    throw new Error('Empty shopping cart')
  }
  try {
    const orders = getOrders(products)
    await submitOrders(Object.values(orders), userProfile)  
    resetCart()
  }catch(e){
    window.alert((e.response && e.response.data && e.response.data.message) || (e.message))
  }
}

export default onClickBuyNow