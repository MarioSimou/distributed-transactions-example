import httpClient from '../../../../utils/httpClient'


const onSubmitForm = (formRef, formValues, onSuccess) => async () => {
  formRef.current.click()

  try { 
    const {data, status} = await httpClient({
      method: 'POST',
      url: new URL('/api/v1/products', process.env.REACT_APP_PRODUCTS_API).toString(),
      data: JSON.stringify({
        name: formValues.productName.value,
        description: formValues.description.value,
        price: formValues.price.value,
        quantity: formValues.quantity.value,
        currency: formValues.currency.value,
        image: formValues.productImage,
      })
    })
    if(status !== 200){
      throw new Error(data.message)
    }
    console.log(data)
    onSuccess()
  }catch(e){
    window.alert((e.response && e.response.data && e.response.data.message) || (e.message))
  }
}


export default onSubmitForm