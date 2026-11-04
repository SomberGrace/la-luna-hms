const FOOD_IMAGES = {
	breakfast:
		'https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?auto=format&fit=crop&w=600&q=80',
	main: 'https://images.unsplash.com/photo-1547496502-affa22d38842?auto=format&fit=crop&w=600&q=80',
	fish: 'https://images.unsplash.com/photo-1519708227418-c8fd9a32b7a2?auto=format&fit=crop&w=600&q=80',
	steak:
		'https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=600&q=80',
	snack:
		'https://images.unsplash.com/photo-1573080496219-bb080dd4f877?auto=format&fit=crop&w=600&q=80',
	drink:
		'https://images.unsplash.com/photo-1613478223719-2ab802602423?auto=format&fit=crop&w=600&q=80',
	beer: 'https://images.unsplash.com/photo-1608270119358-1c4b92b6a50f?auto=format&fit=crop&w=600&q=80',
	cocktail:
		'https://images.unsplash.com/photo-1514362545857-3bc16c4c7d1b?auto=format&fit=crop&w=600&q=80',
	dessert:
		'https://images.unsplash.com/photo-1589367920969-ab8e050bbb04?auto=format&fit=crop&w=600&q=80'
} as const;

const NAMED_FOOD_IMAGES: Record<string, string> = {
	'nsima with stewed chicken': FOOD_IMAGES.main,
	'grilled lake tilapia (chambo)': FOOD_IMAGES.fish,
	'prime beef tenderloin steak': FOOD_IMAGES.steak,
	'full english resort breakfast': FOOD_IMAGES.breakfast,
	'warm mandasi & local honey': FOOD_IMAGES.dessert,
	'crispy chips & peri-peri sauce': FOOD_IMAGES.snack,
	'fresh mango & passion juice': FOOD_IMAGES.drink,
	'coca-cola 500ml':
		'https://images.unsplash.com/photo-1622483767028-3f66f32aef97?auto=format&fit=crop&w=600&q=80',
	'carlsberg green lager (500ml)': FOOD_IMAGES.beer,
	'carlsberg green (500ml)': FOOD_IMAGES.beer,
	'mosi lager (330ml)': FOOD_IMAGES.beer,
	'castle lager (500ml)': FOOD_IMAGES.beer,
	'malawi gin & tonic': FOOD_IMAGES.cocktail,
	'malawi gin 250ml': FOOD_IMAGES.cocktail
};

const CATEGORY_IMAGES: Record<string, string> = {
	breakfast: FOOD_IMAGES.breakfast,
	lunch: FOOD_IMAGES.main,
	dinner: FOOD_IMAGES.steak,
	snack: FOOD_IMAGES.snack,
	beverage: FOOD_IMAGES.drink
};

export function getRestaurantItemImage(name: string, category: string, imageUrl?: string): string {
	const storedImage = imageUrl?.trim();
	if (storedImage) return storedImage;

	return (
		NAMED_FOOD_IMAGES[name.trim().toLowerCase()] ??
		CATEGORY_IMAGES[category.trim().toLowerCase()] ??
		FOOD_IMAGES.main
	);
}

export function getBarItemImage(name: string): string {
	const normalizedName = name.trim().toLowerCase();
	if (
		normalizedName.includes('gin') ||
		normalizedName.includes('spirit') ||
		normalizedName.includes('cocktail')
	) {
		return FOOD_IMAGES.cocktail;
	}
	if (
		normalizedName.includes('beer') ||
		normalizedName.includes('lager') ||
		normalizedName.includes('carlsberg') ||
		normalizedName.includes('mosi') ||
		normalizedName.includes('castle')
	) {
		return FOOD_IMAGES.beer;
	}
	if (normalizedName.includes('cola') || normalizedName.includes('coke')) {
		return NAMED_FOOD_IMAGES['coca-cola 500ml'];
	}
	return FOOD_IMAGES.drink;
}
