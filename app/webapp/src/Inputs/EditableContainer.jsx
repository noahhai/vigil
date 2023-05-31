import React from 'react'
import PropTypes from 'prop-types';
import './EditableContainer.css';

class EditableContainer extends React.Component {
  constructor (props) {
    super(props)

    // init counter
    this.count = 0

    const {initialValue} = this.props;
    // init state
    this.state = {
      edit: false,
      currentValue: initialValue,
    }
  }

  componentWillUnmount () {
    // cancel click callback
    if (this.timeout) clearTimeout(this.timeout)
  }

  handleClick (e) {
    // cancel previous callback
    if (this.timeout) clearTimeout(this.timeout)

    // increment count
    this.count++

    // schedule new callback  [timeBetweenClicks] ms after last click
    this.timeout = setTimeout(() => {
      // listen for double clicks
      if (this.count === 1) {
        // turn on edit mode
        this.setState({
          edit: true,
        })
      }

      // reset count
      this.count = 0
    }, 250) // 250 ms
  }

  handleBlur (e) {
    // handle saving here
    const {onBlur} = this.props;
    const newValue = e.nativeEvent.target.value;
    onBlur(newValue);
    // close edit mode
    this.setState({
      edit: false,
      currentValue: newValue,
    })
  }

  render () {
    const {children, value, label, ...rest} = this.props
    const {edit, currentValue} = this.state

    return (
      <div className="flex-row">
        <div className="flex-items">{label}</div>
        {
          edit ?
          (
            <FieldStyle
              className="flex-items"
              autoFocus
              onBlur={this.handleBlur.bind(this)}
              defaultValue={currentValue}
            />
          ) :
          (
            <div className="flex-items editable"
              onClick={this.handleClick.bind(this)}
            >{currentValue}
            </div>
          )
        }
      </div>
    )
  }
};


EditableContainer.propTypes = {
  initialValue: PropTypes.string,
  label: PropTypes.string.isRequired,
};

EditableContainer.defaultProps = {
  initialValue: '',
};

class FieldStyle extends React.Component {
  componentDidMount () {
    this.ref && this.ref.focus()
  }

  render () {
    const {autoFocus, ...rest} = this.props

    // auto focus
    const ref = autoFocus ? (ref) => { this.ref = ref } : null
    return (
      <input
        ref={ref}
        type="text"
        {...rest}
      />
    )
  }
}


export { EditableContainer};
