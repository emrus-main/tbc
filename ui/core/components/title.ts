import { Component } from '/tbc/core/components/component.js';
import { Spec } from '/tbc/core/proto/common.js';
import {
	classColors,
	getSpecSiteUrl,
	linkedSpecs,
	raidSimIcon,
	raidSimSiteUrl,
	specIconsLarge,
	specNames,
	specToClass,
} from '/tbc/core/proto_utils/utils.js';

declare var tippy: any;

export interface SimLinkOption {
	iconUrl: string,
	href: string,
	text: string,
	color: string,
};

// Dropdown menu for selecting a player.
export class Title extends Component {
	private readonly buttonElem: HTMLElement;
	private readonly dropdownElem: HTMLElement;

  constructor(parent: HTMLElement, currentSpec: Spec | null) {
    super(parent, 'sim-title-root');
		this.rootElem.classList.add('dropdown-root');

    this.rootElem.innerHTML = `
			<div class="dropdown-button sim-title-button"></div>
			<div class="dropdown-panel sim-title-dropdown within-raid-sim-hide"></div>
    `;

		this.buttonElem = this.rootElem.getElementsByClassName('sim-title-button')[0] as HTMLElement;
		this.dropdownElem = this.rootElem.getElementsByClassName('sim-title-dropdown')[0] as HTMLElement;

		this.buttonElem.addEventListener('click', event => {
			event.preventDefault();
		});

		let currentOption = null;
		let otherOptions = [];
		if (currentSpec == null) {
			currentOption = this.makeOptionData(null, true);
			otherOptions = linkedSpecs.map(spec => this.makeOptionData(spec, false));
		} else {
			currentOption = this.makeOptionData(currentSpec, true);
			otherOptions = linkedSpecs
					.filter(spec => spec != currentSpec)
					.map(spec => this.makeOptionData(spec, false))
					.concat([ this.makeOptionData(null, false) ]);
		}

		this.buttonElem.appendChild(Title.makeOptionElem(currentOption));

		const isWithinRaidSim = this.rootElem.closest('.within-raid-sim') != null;
		if (!isWithinRaidSim) {
			otherOptions.forEach(option => this.dropdownElem.appendChild(this.makeOption(option)));
		}
  }

	private makeOptionData(spec: Spec | null, isButton: boolean) {
		if (spec == null) {
			return {
				iconUrl: raidSimIcon,
				href: raidSimSiteUrl,
				text: 'Raid',
				color: isButton ? '' : 'black',
			};
		} else {
			return {
				iconUrl: specIconsLarge[spec],
				href: getSpecSiteUrl(spec),
				text: specNames[spec],
				color: isButton ? '' : classColors[specToClass[spec]],
			};
		}
	}

	private makeOption(data: SimLinkOption): HTMLElement {
		const option = Title.makeOptionElem(data);

		option.addEventListener('click', event => {
			event.preventDefault();
			window.location.href = data.href;
		});

		return option;
	}

	static makeOptionElem(data: SimLinkOption): HTMLElement {
		const option = document.createElement('div');
		option.classList.add('sim-title-option', 'dropdown-option');

		if (data.color) {
			option.style.backgroundColor = data.color;
		}

		const icon = document.createElement('img');
		icon.src = data.iconUrl;
		icon.classList.add('sim-title-icon');
		option.appendChild(icon);

		const labelDiv = document.createElement('div');
		labelDiv.classList.add('sim-title-label-container');
		option.appendChild(labelDiv);

		if (!data.color) { // Hacky check for 'isButton'
			const simLabel = document.createElement('span');
			simLabel.textContent = 'TBC Simulator';
			simLabel.classList.add('sim-title-sim-label', 'sim-title-label');
			labelDiv.appendChild(simLabel);
		}

		const label = document.createElement('span');
		label.textContent = data.text;
		label.classList.add('sim-title-label');
		labelDiv.appendChild(label);

		return option;
	}
}
