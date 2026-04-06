package app

import "fmt"

func (s *Service) UpsertVersion(version string) error {
	target := version
	if target == "" {
		current, err := s.readCurrentVersionName()
		if err != nil {
			return fmt.Errorf("--version not provided and no current version found")
		}
		target = current
	}

	manifest, err := s.captureVersionManifest(target, "upsert version")
	if err != nil {
		return err
	}
	if err := s.refs.SaveVersion(target, manifest); err != nil {
		return err
	}
	return s.setCurrentVersionName(target)
}
