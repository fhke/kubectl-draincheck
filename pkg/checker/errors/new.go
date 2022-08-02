package errors

func For(in error) error {
	if in != nil {
		return &errorImpl{
			error: in,
		}
	} else {
		return nil
	}
}
