package app

func (s *Service) DeleteVersion(version string) error {
	if err := s.refs.DeleteVersion(version); err != nil {
		return err
	}

	st, err := s.state.Get()
	if err != nil {
		return err
	}
	if st.CurrentVersion == version {
		st.CurrentVersion = ""
		if err := s.state.Save(st); err != nil {
			return err
		}
	}
	return nil
}
