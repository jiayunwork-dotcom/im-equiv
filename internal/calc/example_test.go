package calc

import "im-equiv/internal/config"

func example() config.Input {
	return config.Input{P: 4, F: 50, V: 380, Rs: 1.5, Xs: 2.0, Rm: 30.0, Xm: 20.0, R2: 0.8, X2: 1.5, S: 0.03}
}
