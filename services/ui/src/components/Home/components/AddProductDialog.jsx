import React from 'react'
import {
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Typography,
  IconButton,
  InputBase,
  Button,
  TextField,
  FormControl,
  MenuItem,
  Collapse,
  CardMedia,
  CardActionArea,
  Card
} from '@material-ui/core'
import CloseIcon from '@material-ui/icons/Close'
import makeStyles from '@material-ui/styles/makeStyles'
import onSubmitFormProduct from './utils/onSubmitFormProduct.js'
import onChangeFileInput from './utils/onChangeFileInput.js'

const currency = {
  'GBP' : 'GBP',
  'EURO': 'EURO',
  'USD': 'USD'
}

const initProductsValues = {
  productName: {
    touched: false,
    value: ''
  },
  description: {
    touched: false,
    value: '',
  },
  price: {
    touched: false,
    value: 0.0
  },
  quantity: {
    touched: false,
    value: 0,
  },
  currency: {
    touched: false,
    value: currency.GBP
  },
  productImage: ''

}

const AddProductDialog = ({open, onCloseDialog}) => {
  const [formValues, setFormValues] = React.useState(initProductsValues)
  const classes = useStyles()
  const onChangeFieldValue = (key, isNumber = false) => ({target}) => setFormValues({...formValues, [key]: {touched: true, value: isNumber ? +target.value: target.value}}) 
  const formRef = React.useRef()
  const onSuccessSubmission = () => {
    onCloseDialog()
    window.location.reload()
  }
  const onSubmitForm = onSubmitFormProduct(formRef, formValues, onSuccessSubmission)
  const appendField = (key, value) => setFormValues({...formValues, [key]:value})

  return (
    <Dialog open={open} onClose={onCloseDialog} className={classes.root} disableBackdropClick={true}>
      <Paper elevation={3} className={classes.paper}>
        <DialogTitle disableTypography className={classes.dialogTitle} >
          <Typography component="h6" variant="h6">Add Product</Typography>
          <IconButton onClick={onCloseDialog}>
            <CloseIcon/>
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <form className={classes.form} ref={formRef} noValidate={true}>
            <FormControl fullWidth>
              <TextField type="text" 
                         value={formValues.productName.value} 
                         onChange={onChangeFieldValue('productName')} 
                         placeholder="Your Product Name" 
                         label="Product Name" 
                         fullWidth/>
            </FormControl>
            <FormControl fullWidth>
              <TextField multiline 
                        value={formValues.description.value} 
                        onChange={onChangeFieldValue('description')} 
                        placeholder="Your Description" 
                        label="Description"
                        rowsMax={10}
                        fullWidth/>
            </FormControl>
            <FormControl fullWidth>
              <TextField type="number" 
                        value={formValues.price.value} 
                        onChange={onChangeFieldValue('price', true)} 
                        placeholder="Your Price" 
                        label="Price" 
                        fullWidth/>
            </FormControl>
            <FormControl fullWidth>
              <TextField 
                        value={formValues.currency.value} 
                        onChange={onChangeFieldValue('currency')} 
                        placeholder="Your Currency" 
                        label="Currency" 
                        fullWidth
                        select>
                {Object.values(currency).map(curr => {
                  return (<MenuItem key={curr} value={curr}>{curr}</MenuItem>)
                })}
              </TextField>
            </FormControl> 
            <FormControl fullWidth>
              <TextField type="number" 
                        value={formValues.quantity.value} 
                        onChange={onChangeFieldValue('quantity', true)} 
                        placeholder="Your Quantity" 
                        label="Quantity" 
                        fullWidth/>
            </FormControl>
            <FormControl fullWidth>
              <Collapse in={Boolean(formValues.productImage)}>
                {formValues.productImage && 
                  <CardMedia
                  className={classes.productImage}
                  image={formValues.productImage}
                  title="product-image"/>}
              </Collapse>
              <InputBase type="file"
                        id="productImage"
                        onChange={onChangeFileInput(appendField)} 
                        className={classes.hide}/>
              <label htmlFor="productImage">
                <Button variant="contained" color="primary" component="span" fullWidth>Upload image</Button>
              </label>
            </FormControl>  
          </form>
        </DialogContent>
        <DialogActions className={classes.dialogActions}>
          <Button className={classes.button} color="primary" variant="contained" onClick={onSubmitForm} fullWidth>Submit</Button>
          <Button onClick={onCloseDialog} className={classes.button} color="primary" variant="contained" fullWidth>Close</Button>
        </DialogActions>
      </Paper>
    </Dialog>

  )
}

const useStyles = makeStyles(theme => ({
  root: {
  },
  paper: {
    minWidth: 500,
  },
  dialogTitle: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between'
  },
  dialogActions: {
    flexDirection: 'column',
    padding: theme.spacing(3),
  },
  form: {
    display: 'grid',
    gridAutoFlow: 'row',
    gridRowGap: theme.spacing(1),
  },
  button: {
    '&:nth-child(2)': {
      marginTop: theme.spacing(1),
      marginLeft: 0,
    },
  },
  hide: {
    display: 'none'
  },
  productImage: {
    height: 200,
    maxWidth: '100%',
    marginBottom: theme.spacing(1),
  }
}))

export default AddProductDialog