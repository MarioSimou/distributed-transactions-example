import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import classnames from 'classnames'
import onClickDeleteProduct from './utils/onClickDeleteProduct'

const ProductCard = ({id, name:title, description, image, onDeleteProduct}) => {
  const classes = useStyles();
  const [mouseEnter, setMouseEnter] = React.useState(false)
  const onMouseEnter = () => setMouseEnter(true)
  const onMouseLeave = () => setMouseEnter(false)
  
  return (
    <Card className={classes.root} onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      <CardActionArea>
        <CardMedia
          className={classes.media}
          image={image}
          title={title}
        />
        <CardContent className={classnames(classes.cardContent, !mouseEnter && classes.hide)}>
          <Typography gutterBottom variant="h5" component="h2" align="center">
            {title}
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">{description}</Typography>
        </CardContent>
      </CardActionArea>
      <CardActions>
        <Button size="small" color="primary">
          View
        </Button>
        <Button size="small" color="primary" onClick={onClickDeleteProduct(id,onDeleteProduct)}>
          Delete
        </Button>
        <Button size="small" color="primary">
          Add to Cart
        </Button>
      </CardActions>
    </Card>
  );
}

const useStyles = makeStyles(theme => ({
  root: {
    maxWidth: 350,
    position: 'relative',
    [theme.breakpoints.down('sm')]: {
      maxWidth: '100%',
    }
  },
  cardContent: {
    backgroundColor: 'rgba(255,255,255,0.4)',
    width: '100%',
    height: '100%',
    position: 'absolute',
    top: 0,
    left: 0,
    padding: theme.spacing(2),
  },
  hide: {
    display: 'none'
  },
  media: {
    height: 250,
  },
}))

export default ProductCard